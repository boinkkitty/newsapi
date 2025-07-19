//go:build integration

package test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func forward(t testing.TB, serviceNamespace string, serviceName string, containerPort int) string {
	t.Helper()

	// Load client configs
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingClientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	clientConfig, err := loadingClientConfig.ClientConfig()
	require.NoError(t, err)

	// Use spdy transport and upgrader
	transport, upgrader, err := spdy.RoundTripperFor(clientConfig)
	require.NoError(t, err)

	// Provide buffers for errors and outputs
	var bufErr, bufOut bytes.Buffer
	stopChan := make(chan struct{}, 1)
	readChan := make(chan struct{}, 1)

	// Get a free port on local machine
	fp := freePort(t)
	host := strings.TrimPrefix(clientConfig.Host, "https://")
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", serviceNamespace, randomPod(t, clientConfig, serviceNamespace, serviceName).GetName())

	//Create port forwarder
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: host})
	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", fp, containerPort)}, stopChan, readChan, &bufOut, &bufErr)
	require.NoError(t, err)

	go func() {
		err = forwarder.ForwardPorts()
		require.NoError(t, err)
		<-t.Context().Done()
		close(stopChan)
	}()

	<-readChan
	return fmt.Sprintf("http://localhost:%d", fp)
}

func freePort(tb testing.TB) int {
	tb.Helper()

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(tb, err)

	fp := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	return fp
}

func randomPod(tb testing.TB, config *rest.Config, serviceNamespace string, serviceName string) *core.Pod {
	tb.Helper()

	ctx := tb.Context()
	// Get k8s client set
	clientSet, err := kubernetes.NewForConfig(config)
	require.NoError(tb, err)

	// Get the service
	service, err := clientSet.CoreV1().Services(serviceNamespace).Get(ctx, serviceName, meta.GetOptions{})
	require.NoError(tb, err)
	require.NotNil(tb, service)

	labels := make([]string, 0)
	for key, value := range service.Spec.Selector {
		labels = append(labels, key+"="+value)
	}

	// Get all pods matching service labels
	pods, err := clientSet.CoreV1().Pods(serviceNamespace).List(ctx, meta.ListOptions{LabelSelector: strings.Join(labels, " "), Limit: 1})
	require.NoError(tb, err)
	require.NotEmpty(tb, pods.Items)

	// Return first pod
	return &pods.Items[0]
}

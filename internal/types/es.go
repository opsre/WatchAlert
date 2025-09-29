package types

import (
	"crypto/md5"
	"encoding/hex"
	"watchAlert/pkg/tools"
)

type ESQueryFilter struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type ESQueryResponse struct {
	Source source `json:"_source"`
}

type source struct {
	Topic            string `json:"topic"`
	Index            string `json:"index"`
	DockerContainer  string `json:"docker_container"`
	K8sPodNamespace  string `json:"k8s_pod_namespace"`
	K8sPod           string `json:"k8s_pod"`
	K8sContainerName string `json:"k8s_container_name"`
	Message          string `json:"message"`
}

func (r ESQueryResponse) GetMetric() map[string]interface{} {
	return map[string]interface{}{
		"Topic": r.Source.Topic,
		"Index": r.Source.Index,
	}
}

func (r ESQueryResponse) GetFingerprint() string {
	newMetric := map[string]interface{}{
		"Topic": r.Source.Topic,
		"Index": r.Source.Index,
	}
	h := md5.New()
	h.Write([]byte(tools.JsonMarshalToString(newMetric)))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}

func (r ESQueryResponse) GetAnnotations() string {
	s := tools.JsonMarshalToString(r)
	annotations := tools.FormatJson(s)
	return annotations
}

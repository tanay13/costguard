package types

type Provider string

const (
	ProviderKubernetes Provider = "kubernetes"
	ProviderAWSLambda  Provider = "aws_lambda"
	ProviderAWSEC2     Provider = "aws_ec2"
	ProviderVercel     Provider = "vercel"
)

type MetricCollection struct {
	Provider  Provider        `json:"provider"`
	Resource  string          `json:"resource"`
	TimeStamp int64           `json:"timestamp"`
	Metrics   ResourceMetrics `json:"resource_metrics"`
}

type K8sResourceMetrics struct {
	Resource string  `json:"resource,omitempty"`
	CpuMilli float64 `json:"cpu_milli,omitempty"`
	MemoryGB float64 `json:"memory_gb,omitempty"`
}

type LambdaResourceMetrics struct {
	DurationMs  float64 `json:"duration_ms,omitempty"`
	Invocations float64 `json:"invocations,omitempty"`
}

type VMResourceMetrics struct {
	CpuPercent float64 `json:"cpu_percent,omitempty"`
	NetworkGB  float64 `json:"network_gb,omitempty"`
	DiskGB     float64 `json:"disk_gb,omitempty"`
}

type VercelResourceMetrics struct {
	TotalMs    float64 `json:"total_ms,omitempty"`
	ColdStarts float64 `json:"cold_starts,omitempty"`
}

type ResourceMetrics struct {
	K8sResourceMetrics    K8sResourceMetrics    `json:"k8s_resource,omitempty"`
	LambdaResourceMetrics LambdaResourceMetrics `json:"lambda_resource,omitempty"`
	VMResourceMetrics     VMResourceMetrics     `json:"vm_resource,omitempty"`
	VercelResourceMetrics VercelResourceMetrics `json:"vercel_resource,omitempty"`
}

type Requests struct {
	CpuMilli float64 `json:"cpu_milli"`
	MemoryGB float64 `json:"memory_gb"`
}

type MetricStat struct {
	P50 float64 `json:"p50"`
	P95 float64 `json:"p95"`
	Avg float64 `json:"avg"`
}

type AggregatedMetrics struct {
	Provider Provider              `json:"provider"`
	Resource string                `json:"resource"`
	Metrics  map[string]MetricStat `json:"metrics"`

	RequestedCpuMilli float64 `json:"requested_cpu_milli"`
	RequestedMemoryGB float64 `json:"requested_memory_gb"`

	OptimalCpuMilli float64 `json:"optimal_cpu_milli"`
	OptimalMemoryGB float64 `json:"optimal_memory_gb"`

	CostCurrentUSD float64 `json:"cost_current_usd"`
	CostOptimalUSD float64 `json:"cost_optimal_usd"`
	CostSavingsUSD float64 `json:"cost_savings_usd"`

	DataPoints int `json:"data_points"`
}

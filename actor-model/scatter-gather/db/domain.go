package main

type CollectFeatureRequest struct {
	FeatureType string
}

type FeatureResponse struct {
	FeatureType string
	Value       float32
}

type AggregateRequest struct {
	FeatureTypes []string
}

type AggregateResponse struct {
	Results map[string]float32
}

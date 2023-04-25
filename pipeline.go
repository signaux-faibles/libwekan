package libwekan

import "go.mongodb.org/mongo-driver/bson"

type Pipeline bson.A

func (p *Pipeline) AppendStage(stage bson.M) {
	*p = append(*p, stage)
}

func (p *Pipeline) PrependStage(stage bson.M) {
	*p = append(Pipeline{stage}, (*p)...)
}

func (p *Pipeline) AppendPipeline(pipeline Pipeline) {
	*p = append(*p, pipeline...)
}

func (p *Pipeline) PrependPipeline(pipeline Pipeline) {
	*p = append(pipeline, (*p)...)
}

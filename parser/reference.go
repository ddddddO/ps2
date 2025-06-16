package parser

type references struct {
	nodes           map[int]*ASTNode
	referenceNumber int
}

func newReferences() *references {
	return &references{
		nodes:           map[int]*ASTNode{},
		referenceNumber: 1,
	}
}

func (r *references) store(node *ASTNode) bool {
	if node.Type == ASTNodeTypeReference2 {
		return false
	}

	for _, stored := range r.nodes {
		if node.Index == stored.Index {
			return false
		}
	}

	r.nodes[r.referenceNumber] = node
	r.referenceNumber++
	return true
}

func (r *references) getByID(id int) *ASTNode {
	if ref, ok := r.nodes[id]; ok {
		return ref
	}
	return nil
}

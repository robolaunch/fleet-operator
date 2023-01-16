package error

import "fmt"

type NamespaceNotFoundError struct {
	Err               error
	ResourceKind      string
	ResourceName      string
	ResourceNamespace string
}

func (r *NamespaceNotFoundError) Error() string {
	return fmt.Sprintf("cannot found namespace for resource %v %v", r.ResourceKind, r.ResourceName)
}

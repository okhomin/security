package fs

import "context"

type Filer interface {
	List(ctx context.Context)
	Create(ctx context.Context)
	Update(ctx context.Context)
	ReadContent(ctx context.Context)
}

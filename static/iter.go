package static

import "context"

type IterID struct {
	ids []string
}

func NewIterID(ids []string) IterID {
	return IterID{ids: ids}
}

func (i IterID) NextChan(ctx context.Context) chan string {
	c := make(chan string, 0)
	x := 0
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c)
				return
			default:
				if x >= len(i.ids) {
					x = 0
				}
				c <- i.ids[x]

				x++
			}
		}
	}()
	return c
}

package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, s := range stages {
		out = func(in In) Out {
			inner := make(Bi)

			go func() {
				defer close(inner)
				for v := range in {
					select {
					case <-done:
						return
					default:
					}

					select {
					case inner <- v:
					case <-done:
						return
					}
				}
			}()

			return s(inner)
		}(out)
	}

	return out
}

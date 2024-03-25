package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		out := make(Bi)
		defer close(out)
		return out
	}

	out := in
	for _, s := range stages {
		out = func(in In) Out {
			inner := make(Bi)

			go func() {
				defer func() {
					close(inner)
					for range in { // ожидаем завершения предыдущего обработчика стейджа
					}
				}()
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

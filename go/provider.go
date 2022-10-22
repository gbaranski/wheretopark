package wheretopark

type Provider interface {
	Metadata() map[ID]Metadata
	State() map[ID]State
}

func Create() {

}

func Run(provider Provider) {

}

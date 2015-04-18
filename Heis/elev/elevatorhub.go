package elev

type ElevatorHub struct {
	becomeMaster chan bool
}

func NewElevatorHub(becomeMaster chan bool) *ElevatorHub {
	return &ElevatorHub{
		becomeMaster: becomeMaster,
	}
}

func (eh *ElevatorHub) Run() {

}

<?php

namespace BetterScience;

use \Phalcon\Mvc\User\Component;

abstract class Event extends Component {
	protected $type = false;

	protected $data = [];

	public function getType() {
		return $this->type;
	}

	public function getData() {
		return $this->data;
	}
}
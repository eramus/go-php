<?php

namespace BetterScience\Event;

use \Phalcon\Mvc\User\Component;

abstract class Client extends Component {
	private $events = [];

	public function setEvent(\BetterScience\Event $event = null) {
		if (!is_null($event)) {
			$this->events[] = $event;
		}
	}

	public function getEvents() {
		return $this->events;
	}

	public function count() {
		return count($this->events);
	}

	abstract public function fire();
}
<?php

namespace BetterScience;

use \Phalcon\Mvc\User\Component;

class Events extends Component {
	
	private $events = [];

	public function add(\BetterScience\Event $event = null) {
		if (is_null($event) || !$event->getType()) {
			return;
		}
		if (!isset($this->events[$event->getType()])) {
			$this->events[$event->getType()] = [];
		}

		$event->setDi($this->getDi());
		$this->events[$event->getType()][] = $event;
	}

	public function fire() {
		if (function_exists('fastcgi_finish_request')) {
			fastcgi_finish_request();
		}
		if (empty($this->events)) {
			return;
		}

		$clients = [];
		foreach ($this->events as $clientType => $events) {
			$clientTypeClass = '\\BetterScience\\Event\\Client\\'.$clientType;

			$clients[$clientType] = new $clientTypeClass;

			foreach ($events as $event) {
				$clients[$clientType]->setEvent($event);
			}
		}
		foreach ($clients as $client) {
			$client->setDi($this->getDi());
			$client->fire();
		}
	}

}
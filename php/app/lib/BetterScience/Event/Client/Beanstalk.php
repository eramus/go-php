<?php

namespace BetterScience\Event\Client;

class Beanstalk extends \BetterScience\Event\Client {
	public function fire() {
		if (!$this->count()) {
			return;
		}
		foreach ($this->getEvents() as $event) {
			$eventData = $event->getData();
			if (!isset($eventData['tube'])) {
				continue;
			}

			$this->beanstalk->choose($eventData['tube']);
			$this->beanstalk->put(0, 0, 3600, json_encode([
				'request'	=> REQUEST_ID,
				'action'	=> $eventData['action'],
				'data'		=> $eventData['data']
			]));
		}
	}
}
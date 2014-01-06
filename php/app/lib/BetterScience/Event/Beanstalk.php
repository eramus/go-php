<?php

namespace BetterScience\Event;

class Beanstalk extends \BetterScience\Event {
	const EVENT_TYPE = 'beanstalk';

	public function __construct($data = null) {
		$this->type = self::EVENT_TYPE;
		$this->data = $data;
	}
}
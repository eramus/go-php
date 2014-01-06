<?php

namespace BetterScience;

class Noeq {
	private $socket = null;

	private $host = null;

	private $port = null;

	private $token = '';

	public function __construct($host = '0.0.0.0', $port = 4444, $token = '') {
		$this->host = $host;
		$this->port = $port;
		$this->token = $token;

		$this->connect();
		$this->auth();
	}

	public function __destruct() {
	    fclose($this->socket);
	}

	private function connect() {
		$this->socket = fsockopen($this->host, $this->port, $errno, $errstr);
		if (!$this->socket) {
			throw new Exception('unable to connect: '.$errno.' '.$errstr);
		}
	}

	private function auth() {
		if (empty($this->token)) {
			return;
		}
		$len = strlen($this->token);
		if ($len > 255) {
			throw new Exception('token too long');
		}
		fwrite($this->socket, sprintf("\000%c%s", $len, $this->token));
	}

	public function get($num = 1) {
		if ($num > 255) {
			throw new Exception('request too many ids');
		}

	    fwrite($this->socket, pack('C', $num));
	    $len = (8 * $num);
	    $data = '';
	    do {
	        $part = fread($this->socket, $len);
			$len -= strlen($part);
			$data .= $part;
	    }
	    while ($len > 0);

	    $unpacked = unpack('N'.(2 * $num), $data);
	    $values = [];

	    $i = 1;
	    do {
			$values[] = $unpacked[$i] << 32 | $unpacked[($i+1)];
			$i += 2;
			--$num;
	    }
	    while($num > 0);
	    return $values;
	}

	public function getOne() {
		$values = $this->get();
		return $values[0];
	}
}

<?php

class SignupController extends \Phalcon\Mvc\Controller
{
	public function indexAction()
	{
	}

	public function registerAction()
	{
		print '['.microtime(true).'] START REQUEST: '.REQUEST_ID.'<br>';

		for ($i = 0; $i < 2; ++$i) {
			$data = [
				'request'	=> REQUEST_ID,
				'action'	=> 'signup',
				'data'		=> $this->request->getPost()
			];

			$this->beanstalk->put(0, 0, 3600, json_encode($data));
			print '['.microtime(true).'] '.json_encode($data).'<br>';
		}

		print '['.microtime(true).'] WAIT FOR RESPONSES<br>';
		$this->beanstalk->ignore('default');

		for ($i = 0; $i < 2; ++$i) {
			$this->beanstalk->watch('response_'.REQUEST_ID);
			$job = $this->beanstalk->reserve();
			$this->beanstalk->delete($job['id']);

			print '['.microtime(true).'] '.$job['body'].'<br>';
			$data = json_decode($job['body'], true);
			if (is_null($data)) {
				continue;
			}

/*			$this->events->add(new \BetterScience\Event\Beanstalk([
				'tube' => 'after_request',
				'action' => 'lookup',
				'data' => [
					'user' => $data['data']['id']
				]
			]));*/
		}

		print '['.microtime(true).'] FINISHED REQUEST<br>';
	}
}
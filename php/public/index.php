<?php

try {
    //Register an autoloader
    $loader = new \Phalcon\Loader();
    $loader->registerDirs(array(
        '../app/controllers/',
        '../app/models/',
        '../app/lib/beanstalk/src',
        '../app/lib/BetterScience/*'
    ));

    $loader->registerNamespaces([
        'BetterScience' => '../app/lib/BetterScience'
    ]);

    $loader->register();

    $noeq = new \BetterScience\Noeq;
    define('REQUEST_ID', $noeq->getOne());

    //Create a DI
    $di = new Phalcon\DI\FactoryDefault();

    $events = new \BetterScience\Events;
    $events->setDI($di);
/*    $events->add(new \BetterScience\Event\Beanstalk([
        'tube' => 'default',
        'action' => 'request',
        'data' => [
            'page' => $_SERVER['REQUEST_URI']
        ]
    ]));*/

    $di->set('events', $events);

    $bs = new Socket_Beanstalk;
    $bs->connect();
    $di->set('beanstalk', $bs);

    $di->set('voltService', function($view, $di) {
        $volt = new \Phalcon\Mvc\View\Engine\Volt($view, $di);
        $volt->setOptions([
            'compiledPath' => '../app/compiled/',
            'compiledExtension' => '.compiled',
            'compileAlways' => true,
        ]);
        return $volt;
    });

    //Setting up the view component
    $di->set('view', function(){
        $view = new \Phalcon\Mvc\View();
        $view->setViewsDir('../app/views/');
        $view->registerEngines([
            '.phtml' => 'voltService'
        ]);
        return $view;
    });

    // Handle the request
    $application = new \Phalcon\Mvc\Application($di);
    echo $application->handle()->getContent();

    // fire post request events
    $events->fire();
} catch(\Phalcon\Exception $e) {
     echo 'PhalconException: ', $e->getMessage();
}

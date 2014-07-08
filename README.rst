SinaIP Go
===================

sinaip-go is an IP location querying tool. It provides both cmdline querying
interface and HTTP API.

The location data come from `SINA's IP API <http://int.dpool.sina.com.cn/iplookup/iplookup.php>`_.
The whole data of `SINA's IP API` were dupmed to a binary file named
``ip.dat`` and you can download it from `here <https://github.com/ifduyue/sinaip-generator/releases>`_.
Or, If you prefer, I've put generator scripts at `ifduyue/sinaip-generator <https://github.com/ifduyue/sinaip-generator>`_,
you can do the crawling and generating things yourself.

Install
--------

Source
~~~~~~~

You need go installed and ``GOBIN`` in your ``PATH``. Once that is done,
run the command::

    $ go get github.com/ifduyue/sinaip-go
    $ go install github.com/ifduyue/sinaip-go

Usage
-------

::

    sinaip-go -h
    Usage: sinaip-go [globals] <command> [options]

    httpd command: [options] <addr>

    query command: [options] <ip>...

    global flags:
        -cpus=8 Number of CPUs to use
        -ipdat="" Path to ip.dat, will try to get it from env variable "SINAIPDAT" if left empty.

    examples:
      sinaip-go query 1.2.3.4
      sinaip-go httpd 127.0.0.1:8080

-cpus
~~~~~~

Specifies the number of CPUs to be used internally. It defaults to the amount
of CPUs available in the system.

-ipdat
~~~~~~~

Specifies the path of generated binary file. It defaults to a empty string,
and then the value of ENV variable ``SINAIPDAT``.


Copyright and License
----------------------

Copyright (c) 2014 Yue Du - https://github.com/ifduyue

Licensed under AGPLv3, see LICENSE.

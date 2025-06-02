#!/bin/bash
ruby -v | sed 's/ruby //' | sed 's/(.\+//'

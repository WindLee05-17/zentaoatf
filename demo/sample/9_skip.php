#!/usr/bin/env php
<?php
/**
[case]
title=skip
cid=-1
pid=0

[group]
  1. step 1 >> expect 1

[2. group title 3]
  2.1 step >> expect 2.1
  2.2 step >> expect 2.2

[esac]
*/

checkPreCondition() || print("skip\n");
print(">> ...\n");

function checkPreCondition(){}
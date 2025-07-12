<?
$request = explode("/", $_SERVER['REQUEST_URI']);
preg_replace("/\./","",$request);
unset($request[count($request)-1]);
$path = "serviceB/src".implode("/",$request).".php";
if(file_exists($path)){
 require "serviceB/lib/database.php";
 DataBase::init();
 $res = require $path;
 if(isset($res)) http_response_code($res);
}
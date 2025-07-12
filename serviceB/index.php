<?
$request = explode("/", $_SERVER['PHP_SELF']);
preg_replace("/\./","",$request);
unset($request[0],$request[1]);

// if(end($request)=="")unset($request[count($request)-1]);
$path = "./src/".implode("/",$request).".php";
if(file_exists($path)){
 require "./lib/database.php";
 DataBase::init();
 $res = require $path;
 $stdout = fopen("php://stdout","w+");
 fwrite($stdout ,json_encode($_POST)."\n");
 if(isset($res)) http_response_code($res);
}
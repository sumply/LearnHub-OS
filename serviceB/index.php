<?
define("DEVMODE", true);
$request = explode("/", $_SERVER['PHP_SELF']);
preg_replace("/\./","",$request);
unset($request[0],$request[1]);
function Validate(): array{

}
function flushResponce($string):void {
 $stdout = fopen("php://stdout","w+");
 fwrite($stdout ,"respone from self :".$string);
 echo $string;
}
// if(end($request)=="")unset($request[count($request)-1]);
$path = "./src/".implode("/",$request).".php";
if(file_exists($path)){
 require "./lib/database.php";
 DataBase::init();
 $res = require $path;
 if($res != 1) http_response_code($res);
 $stdout = fopen("php://stdout","w+");
 fwrite($stdout ,json_encode($path."\n".file_get_contents('php://input'))."\n");
} else {
 http_response_code(405);
}
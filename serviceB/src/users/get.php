<?

$input = file_get_contents('php://input');
$data = json_decode($input,true);
if(!$data) return 405;
$data = DataBase::select("users","*","id='".$data["id"]."' OR email='".$data["email"]."'")[0];
flushResponce(json_encode($data));
return 200;
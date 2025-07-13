<?
$input = file_get_contents('php://input');
$data = json_decode($input,true);
if(!$data) return 405;
unset($data["group"]);
if($user = DataBase::Select("users","*","email='".$data["email"]."'"))
$res = DataBase::Insert("users",$data);
echo json_encode($res);
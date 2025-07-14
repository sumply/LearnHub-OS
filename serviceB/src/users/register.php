<?
$input = file_get_contents('php://input');
$data = json_decode($input,true);
if(!$data) return 405;
$validate = [
 "id",
 "firstName",
 "secondName",
 "lastName",
 "email",
 "passwordHash",
 "role",
];

$group = $data["group"];
unset($data["group"]);
$res = $user = [];
if(!($user = DataBase::SelectOne("users","*","email='".$data["email"]."'"))){
 $res = DataBase::Insert("users",$data);
}
$group = DataBase::SelectOne("groups","*","name='".$group."'");
$user = DataBase::SelectOne("users","id","email='".$data["email"]."'");
DataBase::Insert("user_groups",["group_id" => $group["id"], "user_id"=>$user["id"]]);
flushResponce( json_encode(array_merge($user, $res)));
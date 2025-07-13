<?
$input = file_get_contents('php://input');
$data = json_decode($input,true);
if(!$data) return 405;
$group = $data["group"];
unset($data["group"]);

if(!($user = DataBase::Select("users","*","email='".$data["email"]."'"))){
 $res = DataBase::Insert("users",$data);
}
if(!$res) $res = [];
$group = DataBase::Select("groups","*","name='".$group."'")[0];
$user = DataBase::Select("users","id","email='".$data["email"]."'")[0];
DataBase::Insert("user_groups",["group_id" => $group["id"], "user_id"=>$user["id"]]);

echo json_encode(array_merge($user, $res));
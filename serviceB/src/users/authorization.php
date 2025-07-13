<? 
$data = json_decode(file_get_contents('php://input'));
if($res = DataBase::Select("users","*", "password='{$data['password']}' AND email='{$data['email']}'")[0])
echo json_encode($res);
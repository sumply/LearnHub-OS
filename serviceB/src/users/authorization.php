<? 
$data = json_decode(file_get_contents('php://input'));
if($res = DataBase::Select("users","*", "passwordhash='{$data['password']}' AND email='{$data['email']}'")[0])
flushResponce (json_encode($res));
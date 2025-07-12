<?
if(!$_GET["id"]) return;
$data = DataBase::select("users","*","id=".$_GET["id"]);
echo json_encode($data);

<?
class DataBase
{
 static $host = 'localhost';
 static $port = '5432';
 static $dbname = 'main';
 static $user = 'pgadmin';
 static $password = 'pgadmin';
 static PDO $connection;
 static function init()
 {
  $dsn = "pgsql:host=".self::$host.";port=".self::$port.";dbname=".self::$dbname;
  self::$connection = new PDO($dsn,self::$user, self::$password, [
   PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
   PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
   PDO::ATTR_EMULATE_PREPARES   => false, // Recommended for better type safety and performance
  ]);
  $stmt = self::$connection->query("SELECT id, name, email FROM users ORDER BY id DESC LIMIT 5");
  $users = $stmt->fetchAll();
 }
 /**
  * undocumented function summary
  *
  * Undocumented function long description
  *
  * @param array | string $fields Description
  * @return type
  * @throws conditon
  **/
 static function select($table, $fields=[], $conditions = "", $ending ="") : array {
  $stmt = self::$connection->prepare("SELECT :fields FROM :table where :conditions :ending");
  if(is_array($fields)) $fields = implode(",",$fields);
  $stmt->execute([
   ":fields" => $fields,
   ":table" => $table,
   ":conditions" => $conditions,
   ":ending" => $ending
  ]);
  return $stmt->fetchAll();
 }
}

<?
class DataBase
{
 static $host = 'localhost';
 static $port = '5432';
 static $dbname = 'main';
 static $user = 'postgres';
 static $password = 'postgres';
 static PDO $connection;
 static function init()
 {
  $dsn = "pgsql:host=" . self::$host . ";port=" . self::$port . ";dbname=" . self::$dbname;
  self::$connection = new PDO($dsn, self::$user, self::$password, [
   PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
   PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
   PDO::ATTR_EMULATE_PREPARES   => false, // Recommended for better type safety and performance
  ]);
  // $stmt = self::$connection->query("SELECT id, name, email FROM users ORDER BY id DESC LIMIT 5");
  // $users = $stmt->fetchAll();
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
 static function Select($table, $fields = [], $conditions = null, $ending = null): array
 {
  if (is_array($fields)) $fields = implode(",", $fields);
  $query = "SELECT $fields FROM $table";
  if ($conditions) $query.= " where $conditions $ending";
  $res = self::$connection->query($query, PDO::FETCH_ASSOC);
  // $stmt = self::$connection->prepare("SELECT :fields FROM :table where :conditions :ending");
  // $stmt->execute([
  //  ":fields" => $fields,
  //  ":table" => $table,
  //  ":conditions" => $conditions,
  //  ":ending" => $ending
  // ]);
  return $res->fetchAll();
 }
 static function Insert(string $table, array $data): Throwable | array
 {
  $fields = implode(", ", array_keys($data));
  $placeholders = implode(", ", array_map(fn($k) => ":$k", array_keys($data)));
  foreach ($data as $k => $v) {
   $insert[":".$k] = $v;
  }
  $sql = "INSERT INTO $table ($fields) VALUES ($placeholders)";
  $stmt = self::$connection->prepare($sql);
  try {
   
   if($stmt->execute($insert)) return $stmt->fetchAll();
  } catch (\Throwable $th) {
   $message = "";
   switch ($th->getCode()) {
    case 23505:
     $message = "duplicate found!";
     break;
    default:

     break;
   }
   return ["error" => $message];
  }
  return [];
 }

 static function Update(string $table, array $data, string $conditions): bool
 {
  $setParts = [];
  foreach ($data as $key => $value) {
   $setParts[] = "$key = :$key";
  }
  $setClause = implode(", ", $setParts);
  $sql = "UPDATE $table SET $setClause WHERE $conditions";
  $stmt = self::$connection->prepare($sql);
  return $stmt->execute($data);
 }

 static function Remove(string $table, string $conditions): bool
 {
  $sql = "DELETE FROM $table WHERE $conditions";
  return self::$connection->exec($sql) !== false;
 }

 static function DirectQuery(string $query, array $params = []): array|bool
 {
  $stmt = self::$connection->prepare($query);
  $stmt->execute($params);
  if (str_starts_with(strtoupper(trim($query)), 'SELECT')) {
   return $stmt->fetchAll();
  }
  return true;
 }
}

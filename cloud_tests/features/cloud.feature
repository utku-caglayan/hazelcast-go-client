Feature: Query fields of JSON values on a map via SQL JSON operators
  As a Hazelcast user, I should be able to operate on maps
  that already has JSON values
  - by first creating a mapping, and then
  - querying them by SQL statements with JSON operators

  Scenario: Query existing JSON values on a map
    Given there are following entries in a random map "m1"
      | int64 | json |
      | 1023 | {"name":"john", "employeeID":1023, "email":"john.doe@mail.com"} |
      | 1025 | {"name":"jane", "employeeID":1025, "email":"jane.doe@mail.com"} |
    When I create a mapping for "m1"
    """
    CREATE OR REPLACE MAPPING "%s"
    TYPE IMap OPTIONS('keyFormat'='bigint', 'valueFormat'='json');
    """
    And I execute statement for "m1" with result "r1"
    """
    SELECT JSON_QUERY(this, '$')
    FROM "%s" ORDER BY __key DESC
    """
    Then "r1" should be
    |json|
    |{"name":"jane","employeeID":1025,"email":"jane.doe@mail.com"}|
    |{"name":"john","employeeID":1023,"email":"john.doe@mail.com"}|
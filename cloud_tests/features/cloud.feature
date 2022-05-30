Feature: Query fields of JSON values on a map via SQL JSON operators
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
    SELECT JSON_QUERY(this, '$.name')
    FROM "%s" ORDER BY __key DESC
    """
    Then "r1" should be
    |json|
    |"jane"|
    |"john"|
  Scenario: Querying existing JSON values on a map without creating a mapping should fail
    Given there are following entries in a random map "m1"
      | string | json |
      | should | {"name":"fail", "employeeID":1023, "email":"john.doe@mail.com"} |
    When I execute statement for "m1" with result "r1"
    """
    SELECT JSON_QUERY(this, '$')
    FROM "%s" ORDER BY __key DESC
    """
    Then "r1" should be a SQL error
package helper

import (
	"fmt"
	"log"

	"github.com/gocql/gocql" // Scylla Drivers
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* Helper Methods */

/* Creates an index as needed for update and delete methods
*
 * Parameters:
 *   - session: A session for executing select query
 *   - keyspace: The keyspace of the target table
 *   - tableName: The name of the target table
 *   - column: The name of the column for the type
 *
 * Returns:
 *	 - error: returns error. No string response needed
*/
func CreateIndex(session *gocql.Session, keyspace, tableName, column string) error {
	indexQuery := fmt.Sprintf("CREATE INDEX IF NOT EXISTS ON %s.%s (%s)", keyspace, tableName, column)
	if err := session.Query(indexQuery).Exec(); err != nil {
		log.Printf("Error creating index: %s\n query: %s", err, indexQuery)
		return err
	}
	return nil
}

/* Executes selection query for update and delete methods
*
 */
func SelectionQuery(columnType string, ks string, tableName string, column string, constraint string) string {
	var returnQuery string

	switch columnType {
	case "text", "blob", "boolean", "varchar":
		// Selects all the id's that meet the condition
		returnQuery = fmt.Sprintf("SELECT serial_msg FROM %s.%s WHERE %s = '%s'", ks, tableName, column, constraint)
	case "int", "bigint", "float", "double", "uuid":
		returnQuery = fmt.Sprintf("SELECT serial_msg FROM %s.%s WHERE %s = %s", ks, tableName, column, constraint)
	}
	return returnQuery
}

/*
 * getColumnType retrieves the type of a specified column in a Cassandra table
 * by querying the system_schema.columns table
 *
 */
func GetColumnType(session *gocql.Session, keyspace, tableName, columnName string) (string, error) {
	// Selects all the columns and associated types
	iter := session.Query("SELECT column_name, type FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ?", keyspace, tableName).Iter()

	// Loops through returned column names to look for a matching one
	var fetchedColumnName, columnType string
	for iter.Scan(&fetchedColumnName, &columnType) {
		if fetchedColumnName == columnName {
			return columnType, nil
		}
	}

	// If nothing returned or an error is returned return an error
	if err := iter.Close(); err != nil {
		return "", err
	}

	// If didn't error but nothing found returns error
	return "", fmt.Errorf("Column %s not found in table %s", columnName, tableName)
}

/*
 * convertKindAndValueToCQL maps protobuf.Kind and protobuf.Value to CQL data types.
 *
 * It takes a protobuf field type and value, and returns the corresponding CQL data type
 * as a string and the converted value as an interface{}
 *
 * Parameters:
 *   - fieldType: protoreflect.Kind, the protobuf field type.
 *   - value: protoreflect.Value, the value of the protobuf field.
 *
 * Returns:
 *   - string: The CQL data type.
 *   - interface{}: The converted value.
 *
 * TODO: Handle Unknown or other types.
 */
func ConvertKindAndValueToCQL(fieldType protoreflect.Kind, value protoreflect.Value) (string, interface{}) {
	switch fieldType {
	case protoreflect.BoolKind:
		return "BOOLEAN", value.Bool()
	case protoreflect.EnumKind, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "INT", value.Int()
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "BIGINT", value.Int()
	case protoreflect.FloatKind:
		return "FLOAT", value.Float()
	case protoreflect.DoubleKind:
		return "DOUBLE", value.Float()
	case protoreflect.StringKind:
		return "TEXT", value.String()
	case protoreflect.BytesKind:
		return "BLOB", value.Bytes()
	default:
		return "UNKNOWN", nil
	}
}

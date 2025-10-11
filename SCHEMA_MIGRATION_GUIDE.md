# Schema Migration Guide: Different Database Structures → Aiven

## 🎯 Overview

**Situation:**
- Source MySQL has a different schema than your Go application expects
- Need to transform/map old schema to new schema
- Fresh start in Aiven (don't need existing data)
- Go application uses GORM models to define new schema

**Strategy:**
1. Create fresh schema in Aiven using GORM models
2. Map old schema columns → new schema columns
3. Transform and import data with mapping

---

## 📋 Step-by-Step Migration

### **Phase 1: Understand Source Schema** (15 minutes)

#### **Step 1.1: Provide Source Database Details**

Fill in your source database credentials:
```bash
# Save as source_db_config.env
SOURCE_DB_HOST="___________________________"
SOURCE_DB_PORT="___________________________"  # Usually 3306
SOURCE_DB_USER="___________________________"
SOURCE_DB_PASSWORD="___________________________"
SOURCE_DB_NAME="___________________________"
```

---

#### **Step 1.2: Export Source Schema Structure**
```bash
# Create working directory
mkdir -p ~/Desktop/schema_migration_$(date +%Y%m%d)
cd ~/Desktop/schema_migration_*

# Export schema only (no data)
mysqldump -h $SOURCE_DB_HOST \
  -P $SOURCE_DB_PORT \
  -u $SOURCE_DB_USER \
  -p"$SOURCE_DB_PASSWORD" \
  --no-data \
  --skip-triggers \
  $SOURCE_DB_NAME \
  > source_schema.sql

# View schema
cat source_schema.sql
```

---

#### **Step 1.3: List All Tables and Columns**
```bash
# Get detailed table structure
mysql -h $SOURCE_DB_HOST \
  -P $SOURCE_DB_PORT \
  -u $SOURCE_DB_USER \
  -p"$SOURCE_DB_PASSWORD" \
  $SOURCE_DB_NAME \
  -e "
  SELECT 
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_KEY,
    COLUMN_DEFAULT,
    EXTRA
  FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = '$SOURCE_DB_NAME'
  ORDER BY TABLE_NAME, ORDINAL_POSITION;
  " > source_columns_detailed.txt

cat source_columns_detailed.txt
```

**📋 Save this information** - we'll use it to create mapping!

---

### **Phase 2: Define Target Schema in Aiven** (10 minutes)

#### **Step 2.1: Clean Aiven Database**
```bash
# Drop all existing tables in Aiven (fresh start)
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "SET FOREIGN_KEY_CHECKS = 0;
      SET GROUP_CONCAT_MAX_LEN=32768;
      SET @tables = NULL;
      SELECT GROUP_CONCAT(table_name) INTO @tables
        FROM information_schema.tables
        WHERE table_schema = 'defaultdb';
      SELECT IFNULL(@tables,'dummy') INTO @tables;
      SET @tables = CONCAT('DROP TABLE IF EXISTS ', @tables);
      PREPARE stmt FROM @tables;
      EXECUTE stmt;
      DEALLOCATE PREPARE stmt;
      SET FOREIGN_KEY_CHECKS = 1;"

# Verify empty
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "SHOW TABLES;"
```

**✅ Expected:** Should show no tables

---

#### **Step 2.2: Create New Schema from GORM Models**
```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go

# Run GORM auto-migration to create tables based on Go models
DB_USER=avnadmin \
DB_PASSWORD=AVNS_7QDxgZDlRhQXAx3QV4z \
DB_HOST=nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
DB_PORT=22013 \
DB_NAME=defaultdb \
ENVIRONMENT=development \
go run migrate/migrate.go

# Should see:
# ✅ Database connected successfully
# ✅ Migration completed successfully
```

---

#### **Step 2.3: Export New Aiven Schema for Comparison**
```bash
cd ~/Desktop/schema_migration_*

# Export Aiven schema
mysqldump -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  --no-data \
  defaultdb \
  > target_schema.sql

# Get Aiven table structure
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "
  SELECT 
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE
  FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = 'defaultdb'
  ORDER BY TABLE_NAME, ORDINAL_POSITION;
  " > target_columns.txt

cat target_columns.txt
```

---

### **Phase 3: Create Schema Mapping** (30 minutes)

#### **Step 3.1: Compare Schemas Side-by-Side**

Create a mapping file for each table:

```bash
cd ~/Desktop/schema_migration_*

cat > schema_mapping.md << 'EOF'
# Schema Mapping: Source → Target (Aiven)

## Table: users → user_auths

| Source Column | Source Type | Target Column | Target Type | Transformation |
|---------------|-------------|---------------|-------------|----------------|
| [old_col_1]   | [type]      | id            | BIGINT      | AUTO_INCREMENT |
| [old_col_2]   | [type]      | first_name    | VARCHAR     | TRIM, LOWER    |
| [old_col_3]   | [type]      | last_name     | VARCHAR     | TRIM           |
| [old_col_4]   | [type]      | email         | VARCHAR     | TRIM, LOWER    |
| [old_col_5]   | [type]      | password      | VARCHAR     | ALREADY HASHED |
| [old_col_6]   | [type]      | user_type     | VARCHAR     | MAP: 1→ADMIN, 2→CLIENT |

**Notes:**
- Drop columns: [list old columns you don't need]
- New columns with defaults: token='', refresh_token=''
- created_at, updated_at will be set to NOW()

---

## Table: [old_clients_table] → clients

| Source Column | Source Type | Target Column | Target Type | Transformation |
|---------------|-------------|---------------|-------------|----------------|
| ...           | ...         | ...           | ...         | ...            |

---

## Table: [old_recipes_table] → recipes

...

EOF

# Open in editor to fill in
open schema_mapping.md
# OR: nano schema_mapping.md
```

**📝 For each old table, document:**
1. What's the old table name?
2. What's the new table name in Aiven?
3. How do columns map (old → new)?
4. What transformations are needed?
5. What columns to drop?
6. What new columns need default values?

---

### **Phase 4: Create Migration SQL Script** (45 minutes)

Once you have the mapping, create a transformation script:

```bash
cd ~/Desktop/schema_migration_*

cat > migrate_data.sql << 'EOF'
-- Migration Script: Source → Aiven
-- Created: $(date)

-- ================================================================
-- TABLE 1: user_auths
-- ================================================================

-- Insert users with column mapping
INSERT INTO defaultdb.user_auths (
    first_name,
    last_name,
    email,
    password,
    user_type,
    token,
    refresh_token,
    created_at,
    updated_at
)
SELECT 
    TRIM([source_firstname_col]) as first_name,
    TRIM([source_lastname_col]) as last_name,
    LOWER(TRIM([source_email_col])) as email,
    [source_password_col] as password,
    CASE [source_usertype_col]
        WHEN 1 THEN 'ADMIN'
        WHEN 2 THEN 'CLIENT'
        ELSE 'CLIENT'
    END as user_type,
    '' as token,
    '' as refresh_token,
    COALESCE([source_created_col], NOW()) as created_at,
    COALESCE([source_updated_col], NOW()) as updated_at
FROM [SOURCE_DATABASE].[SOURCE_USER_TABLE]
WHERE [source_email_col] IS NOT NULL;  -- Only migrate valid users

-- ================================================================
-- TABLE 2: clients
-- ================================================================

INSERT INTO defaultdb.clients (
    name,
    email,
    age,
    phone_number,
    city,
    height,
    starting_weight,
    is_active,
    date_of_joining,
    created_at,
    updated_at
)
SELECT 
    CONCAT(TRIM([source_first_name]), ' ', TRIM([source_last_name])) as name,
    LOWER(TRIM([source_email])) as email,
    [source_age] as age,
    [source_phone] as phone_number,
    [source_city] as city,
    [source_height] as height,
    [source_weight] as starting_weight,
    COALESCE([source_is_active], 1) as is_active,
    COALESCE([source_join_date], NOW()) as date_of_joining,
    COALESCE([source_created], NOW()) as created_at,
    COALESCE([source_updated], NOW()) as updated_at
FROM [SOURCE_DATABASE].[SOURCE_CLIENT_TABLE];

-- ================================================================
-- TABLE 3: recipes
-- ================================================================

INSERT INTO defaultdb.recipes (
    name,
    image_url,
    created_at,
    updated_at
)
SELECT 
    TRIM([source_recipe_name]) as name,
    [source_image_path] as image_url,
    COALESCE([source_created], NOW()) as created_at,
    COALESCE([source_updated], NOW()) as updated_at
FROM [SOURCE_DATABASE].[SOURCE_RECIPE_TABLE];

-- ================================================================
-- Add more tables as needed...
-- ================================================================

EOF
```

**✍️ You need to:**
1. Replace `[source_*_col]` with actual source column names from `source_columns_detailed.txt`
2. Replace `[SOURCE_DATABASE]` with your source database name
3. Replace `[SOURCE_*_TABLE]` with actual source table names
4. Add any custom transformations (CASE statements, CONCAT, etc.)

---

### **Phase 5: Execute Migration** (20 minutes)

#### **Step 5.1: Backup Source Data First**
```bash
cd ~/Desktop/schema_migration_*

# Full backup of source (with data)
mysqldump -h $SOURCE_DB_HOST \
  -P $SOURCE_DB_PORT \
  -u $SOURCE_DB_USER \
  -p"$SOURCE_DB_PASSWORD" \
  --single-transaction \
  $SOURCE_DB_NAME \
  > source_full_backup_$(date +%Y%m%d_%H%M%S).sql

gzip source_full_backup_*.sql
```

---

#### **Step 5.2: Test Migration Script Syntax**
```bash
# Test syntax on source database (dry run)
mysql -h $SOURCE_DB_HOST \
  -P $SOURCE_DB_PORT \
  -u $SOURCE_DB_USER \
  -p"$SOURCE_DB_PASSWORD" \
  --show-warnings \
  < migrate_data.sql > migration_test.log 2>&1

# Check for errors
cat migration_test.log | grep -i "error"
```

---

#### **Step 5.3: Execute Migration (Source → Aiven)**

**Option A: Direct Cross-Database INSERT** (if source is accessible from your machine)
```bash
# Run migration script
# This will read from source and insert into Aiven
mysql -h $SOURCE_DB_HOST \
  -P $SOURCE_DB_PORT \
  -u $SOURCE_DB_USER \
  -p"$SOURCE_DB_PASSWORD" \
  < migrate_data.sql

# Check progress (in another terminal)
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "SELECT 'user_auths' as tbl, COUNT(*) as cnt FROM user_auths
      UNION ALL SELECT 'clients', COUNT(*) FROM clients
      UNION ALL SELECT 'recipes', COUNT(*) FROM recipes;"
```

**Option B: Export → Transform → Import** (more controlled)
```bash
# 1. Export source data as CSV
mysql -h $SOURCE_DB_HOST \
  -P $SOURCE_DB_PORT \
  -u $SOURCE_DB_USER \
  -p"$SOURCE_DB_PASSWORD" \
  $SOURCE_DB_NAME \
  -e "SELECT * FROM [source_users_table]" \
  > source_users.csv

# 2. Transform with script (create Python/Go script if needed)
# See Phase 6 for transformation script example

# 3. Import to Aiven
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  --local-infile=1 \
  defaultdb \
  -e "LOAD DATA LOCAL INFILE 'transformed_users.csv'
      INTO TABLE user_auths
      FIELDS TERMINATED BY ','
      ENCLOSED BY '\"'
      LINES TERMINATED BY '\n'
      IGNORE 1 ROWS;"
```

---

### **Phase 6: Data Transformation Script** (Optional, if SQL too complex)

If transformations are complex, create a Go script:

```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go

cat > migrate/transform_and_import.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"
    "strings"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Source database connection
    sourceDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        "SOURCE_USER", "SOURCE_PASS", "SOURCE_HOST", "SOURCE_PORT", "SOURCE_DB")
    sourceDB, err := sql.Open("mysql", sourceDSN)
    if err != nil {
        log.Fatal("Source connection failed:", err)
    }
    defer sourceDB.Close()

    // Target (Aiven) connection
    targetDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=skip-verify&parseTime=true",
        "avnadmin", "AVNS_7QDxgZDlRhQXAx3QV4z", 
        "nutriediet-mysql-ishitagupta-5564.f.aivencloud.com", 
        "22013", "defaultdb")
    targetDB, err := sql.Open("mysql", targetDSN)
    if err != nil {
        log.Fatal("Target connection failed:", err)
    }
    defer targetDB.Close()

    // Example: Migrate users
    migrateUsers(sourceDB, targetDB)
    
    // Add more migration functions...
    // migrateClients(sourceDB, targetDB)
    // migrateRecipes(sourceDB, targetDB)
    
    log.Println("✅ Migration completed!")
}

func migrateUsers(source, target *sql.DB) {
    log.Println("Migrating users...")
    
    // Query source
    rows, err := source.Query(`
        SELECT [old_col1], [old_col2], [old_col3], [old_col4]
        FROM [old_users_table]
    `)
    if err != nil {
        log.Fatal("Query failed:", err)
    }
    defer rows.Close()

    // Prepare insert statement
    stmt, err := target.Prepare(`
        INSERT INTO user_auths (first_name, last_name, email, password, user_type, token, refresh_token, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, '', '', NOW(), NOW())
    `)
    if err != nil {
        log.Fatal("Prepare failed:", err)
    }
    defer stmt.Close()

    count := 0
    for rows.Next() {
        var oldCol1, oldCol2, oldCol3, oldCol4 string
        if err := rows.Scan(&oldCol1, &oldCol2, &oldCol3, &oldCol4); err != nil {
            log.Println("Scan error:", err)
            continue
        }

        // Transform data
        firstName := strings.TrimSpace(oldCol1)
        lastName := strings.TrimSpace(oldCol2)
        email := strings.ToLower(strings.TrimSpace(oldCol3))
        password := oldCol4
        userType := "CLIENT" // Transform as needed

        // Insert to target
        _, err := stmt.Exec(firstName, lastName, email, password, userType)
        if err != nil {
            log.Printf("Insert failed for %s: %v\n", email, err)
            continue
        }
        count++
    }

    log.Printf("✅ Migrated %d users\n", count)
}

// Add similar functions for other tables...
EOF

# Run the transformation script
go run migrate/transform_and_import.go
```

---

### **Phase 7: Verification** (15 minutes)

#### **Step 7.1: Compare Record Counts**
```bash
# Source counts
mysql -h $SOURCE_DB_HOST \
  -P $SOURCE_DB_PORT \
  -u $SOURCE_DB_USER \
  -p"$SOURCE_DB_PASSWORD" \
  $SOURCE_DB_NAME \
  -e "
  SELECT '[source_users_table]' as table_name, COUNT(*) as count FROM [source_users_table]
  UNION ALL SELECT '[source_clients_table]', COUNT(*) FROM [source_clients_table]
  UNION ALL SELECT '[source_recipes_table]', COUNT(*) FROM [source_recipes_table];
  " > source_final_counts.txt

# Aiven counts
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "
  SELECT 'user_auths' as table_name, COUNT(*) as count FROM user_auths
  UNION ALL SELECT 'clients', COUNT(*) FROM clients
  UNION ALL SELECT 'recipes', COUNT(*) FROM recipes;
  " > aiven_final_counts.txt

# Compare
echo "=== SOURCE COUNTS ==="
cat source_final_counts.txt
echo ""
echo "=== AIVEN COUNTS ==="
cat aiven_final_counts.txt
```

**✅ Verify:** Counts should match (or be close if you filtered some data)

---

#### **Step 7.2: Spot Check Data Quality**
```bash
# Check a few random records
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "
  SELECT * FROM user_auths LIMIT 5;
  SELECT * FROM clients LIMIT 5;
  SELECT * FROM recipes LIMIT 5;
  "
```

**✅ Check for:**
- [ ] Names look correct
- [ ] Emails are lowercase and trimmed
- [ ] No NULL values where NOT NULL required
- [ ] Dates are reasonable
- [ ] Foreign keys will work (IDs match)

---

### **Phase 8: Test Application** (10 minutes)

```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go

# Make sure .env points to Aiven
cat .env | grep DB_

# Start application
go run main.go

# Test in browser or with curl
curl http://localhost:8080/health

# Test login with migrated user
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**✅ Testing:**
- [ ] Application starts
- [ ] Can query database
- [ ] Login works
- [ ] Can view clients
- [ ] Can create/edit records

---

## 📋 Quick Reference: Common Transformations

### **String Transformations**
```sql
-- Trim whitespace
TRIM(column_name)

-- Lowercase
LOWER(column_name)

-- Uppercase
UPPER(column_name)

-- Concatenate
CONCAT(first_name, ' ', last_name)

-- Substring
SUBSTRING(column_name, 1, 50)

-- Replace
REPLACE(column_name, 'old', 'new')
```

### **Type Conversions**
```sql
-- String to INT
CAST(column_name AS UNSIGNED)

-- String to DATE
STR_TO_DATE(column_name, '%Y-%m-%d')

-- NULL handling
COALESCE(column_name, default_value)

-- Boolean conversion
CASE WHEN column_name = '1' THEN true ELSE false END
```

### **Conditional Mapping**
```sql
-- Map values
CASE column_name
    WHEN 'old_value1' THEN 'new_value1'
    WHEN 'old_value2' THEN 'new_value2'
    ELSE 'default_value'
END
```

---

## 🆘 What I Need From You

To help you create the actual migration script, please provide:

### **1. Source Database Connection Details**
```
Host: ____________________
Port: ____________________
User: ____________________
Password: ________________
Database: ________________
```

### **2. Source Table Names**
What are your current table names? Example:
- tbl_users?
- user_table?
- members?
etc.

### **3. Key Column Differences**
For each important table, what are the old column names?

**Example:**
```
Old table: tbl_users
- user_id (INT)
- fname (VARCHAR)
- lname (VARCHAR)
- email_address (VARCHAR)
- pwd_hash (VARCHAR)
- account_type (INT) -- 1=admin, 2=client
- created (DATETIME)

New table: user_auths
- id (BIGINT)
- first_name (VARCHAR)
- last_name (VARCHAR)
- email (VARCHAR)
- password (VARCHAR)
- user_type (VARCHAR) -- 'ADMIN' or 'CLIENT'
- created_at (DATETIME)
```

Once you provide this information, I can create the **exact SQL migration script** tailored to your schema! 🚀

---

**Created:** October 5, 2025  
**Version:** 1.0  
**Status:** Template - needs customization


# Database Migration Plan: External MySQL → Aiven Cloud Database

## 📋 Overview

**Migration Direction:** External MySQL → Aiven MySQL
- **Target (Aiven)**: `nutriediet-mysql-ishitagupta-5564.f.aivencloud.com:22013`
- **Source**: Your existing MySQL database (to be configured)
- **Goal**: Migrate all data from your existing MySQL to Aiven cloud database

---

## 🎯 Step-by-Step Migration Plan

### **Phase 1: Information Gathering** (10 minutes)

#### **Step 1.1: Get Source Database Details**
We need the following information about your source MySQL database:

```
Source Database Credentials:
- Host/IP: _________________________
- Port: ___________________________ (default: 3306)
- Username: _______________________
- Password: _______________________
- Database Name: __________________
```

#### **Step 1.2: Test Source Database Connection**
```bash
# Replace with your actual credentials
mysql -h YOUR_SOURCE_HOST \
  -P YOUR_SOURCE_PORT \
  -u YOUR_SOURCE_USER \
  -p'YOUR_SOURCE_PASSWORD' \
  -e "SHOW DATABASES;"
```

**✅ Checklist:**
- [ ] Can connect to source database
- [ ] Can list databases
- [ ] Have proper read permissions

---

### **Phase 2: Backup & Analyze Source** (20 minutes)

#### **Step 2.1: Create Full Backup of Source Database**
```bash
# Create backup directory
mkdir -p ~/Desktop/nutriediet_migration_$(date +%Y%m%d_%H%M%S)
cd ~/Desktop/nutriediet_migration_*

# Full backup (schema + data)
mysqldump -h YOUR_SOURCE_HOST \
  -P YOUR_SOURCE_PORT \
  -u YOUR_SOURCE_USER \
  -p'YOUR_SOURCE_PASSWORD' \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  --databases YOUR_SOURCE_DATABASE_NAME \
  > source_full_backup.sql

# Compress backup
gzip source_full_backup.sql

# Verify backup
ls -lh source_full_backup.sql.gz
```

**✅ Verification:**
- [ ] Backup file created
- [ ] Backup size is reasonable (>100KB)
- [ ] No errors during dump

---

#### **Step 2.2: Analyze Source Database Structure**
```bash
# Get table list
mysql -h YOUR_SOURCE_HOST \
  -P YOUR_SOURCE_PORT \
  -u YOUR_SOURCE_USER \
  -p'YOUR_SOURCE_PASSWORD' \
  YOUR_SOURCE_DATABASE_NAME \
  -e "SHOW TABLES;" > source_tables.txt

cat source_tables.txt

# Get record counts
mysql -h YOUR_SOURCE_HOST \
  -P YOUR_SOURCE_PORT \
  -u YOUR_SOURCE_USER \
  -p'YOUR_SOURCE_PASSWORD' \
  YOUR_SOURCE_DATABASE_NAME \
  -e "
  SELECT 
    table_name,
    table_rows,
    ROUND((data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
  FROM information_schema.tables
  WHERE table_schema = 'YOUR_SOURCE_DATABASE_NAME'
  ORDER BY table_rows DESC;
  " > source_statistics.txt

cat source_statistics.txt
```

**✅ Checklist:**
- [ ] Table list documented
- [ ] Record counts recorded
- [ ] Database size calculated

---

### **Phase 3: Prepare Aiven Target** (15 minutes)

#### **Step 3.1: Test Aiven Connection**
```bash
# Test connection to Aiven
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  -e "SELECT 'Aiven connection successful!' as status;"
```

**✅ Verification:**
- [ ] Connected to Aiven successfully
- [ ] SSL working properly

---

#### **Step 3.2: Check Current Aiven Database State**
```bash
# List existing tables
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "SHOW TABLES;" > aiven_current_tables.txt

cat aiven_current_tables.txt

# Get current record counts (if tables exist)
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "
  SELECT 
    table_name,
    table_rows
  FROM information_schema.tables
  WHERE table_schema = 'defaultdb'
  ORDER BY table_rows DESC;
  " > aiven_current_counts.txt

cat aiven_current_counts.txt
```

**⚠️ Important Decision Point:**

If Aiven already has tables/data:
- **Option A**: Drop existing tables and import fresh (clean slate)
- **Option B**: Merge data (keep existing + add new) - more complex
- **Option C**: Backup current Aiven data first, then replace

**Recommended**: Option A (fresh import) unless you need existing Aiven data

---

#### **Step 3.3: Backup Current Aiven Data (Optional but Recommended)**
```bash
# Backup current Aiven state before migration
mysqldump -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  --databases defaultdb \
  > aiven_before_migration.sql

gzip aiven_before_migration.sql
```

---

### **Phase 4: Data Migration** (30 minutes)

#### **Step 4.1: Prepare Migration File**
```bash
cd ~/Desktop/nutriediet_migration_*

# Extract backup
gunzip source_full_backup.sql.gz

# Optional: Edit database name if needed
# If source database name differs from 'defaultdb', replace it:
sed -i.bak 's/USE `YOUR_OLD_DB_NAME`/USE `defaultdb`/g' source_full_backup.sql
sed -i.bak 's/CREATE DATABASE.*YOUR_OLD_DB_NAME/CREATE DATABASE IF NOT EXISTS `defaultdb`/g' source_full_backup.sql
```

---

#### **Step 4.2: Clear Aiven Database (Clean Slate Approach)**
```bash
# Get list of all tables to drop
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

# Verify tables dropped
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "SHOW TABLES;"
```

**✅ Expected**: No tables should be listed

---

#### **Step 4.3: Import Data to Aiven**
```bash
# Import full backup to Aiven
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  < source_full_backup.sql

# This may take 5-30 minutes depending on data size
# Watch for any errors
```

**⚠️ Common Errors & Solutions:**

**Error**: `Unknown database 'YOUR_OLD_DB_NAME'`
**Solution**: Edit SQL file to change database name to `defaultdb`

**Error**: `Access denied` or permission errors
**Solution**: Make sure using `avnadmin` user with full privileges

**Error**: `Max allowed packet` errors
**Solution**: Import in smaller chunks (see alternative method below)

---

#### **Step 4.4: Alternative - Import Table by Table (If Full Import Fails)**
```bash
# Export each table separately from source
for table in $(mysql -h YOUR_SOURCE_HOST -P YOUR_SOURCE_PORT -u YOUR_SOURCE_USER -p'YOUR_SOURCE_PASSWORD' YOUR_SOURCE_DATABASE_NAME -e "SHOW TABLES" -s --skip-column-names); do
  echo "Exporting table: $table"
  mysqldump -h YOUR_SOURCE_HOST \
    -P YOUR_SOURCE_PORT \
    -u YOUR_SOURCE_USER \
    -p'YOUR_SOURCE_PASSWORD' \
    --single-transaction \
    YOUR_SOURCE_DATABASE_NAME $table > ${table}.sql
  
  echo "Importing table: $table to Aiven"
  mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
    -P 22013 \
    -u avnadmin \
    -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
    --ssl-mode=REQUIRED \
    defaultdb < ${table}.sql
done
```

---

### **Phase 5: Verification** (15 minutes)

#### **Step 5.1: Verify Table Structure**
```bash
# List tables in Aiven
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "SHOW TABLES;" > aiven_after_migration_tables.txt

# Compare with source
echo "=== SOURCE TABLES ==="
cat source_tables.txt
echo ""
echo "=== AIVEN TABLES (AFTER MIGRATION) ==="
cat aiven_after_migration_tables.txt

# Check for differences
diff source_tables.txt aiven_after_migration_tables.txt
```

**✅ Expected**: All tables should match

---

#### **Step 5.2: Verify Record Counts**
```bash
# Get Aiven record counts
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "
  SELECT 
    table_name,
    table_rows
  FROM information_schema.tables
  WHERE table_schema = 'defaultdb'
  ORDER BY table_name;
  " > aiven_after_migration_counts.txt

# Compare
echo "=== SOURCE COUNTS ==="
cat source_statistics.txt
echo ""
echo "=== AIVEN COUNTS (AFTER MIGRATION) ==="
cat aiven_after_migration_counts.txt
```

**✅ Verification:**
- [ ] All tables present
- [ ] Record counts match (or close)
- [ ] No negative numbers

---

#### **Step 5.3: Spot Check Sample Data**
```bash
# Check a few records from key tables
# Replace TABLE_NAME with your actual table names

mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb \
  -e "
  -- Check user_auths (if exists)
  SELECT COUNT(*) as user_count FROM user_auths;
  SELECT * FROM user_auths LIMIT 3;
  
  -- Check clients (if exists)
  SELECT COUNT(*) as client_count FROM clients;
  SELECT * FROM clients LIMIT 3;
  
  -- Check recipes (if exists)
  SELECT COUNT(*) as recipe_count FROM recipes;
  SELECT * FROM recipes LIMIT 3;
  "
```

**✅ Checklist:**
- [ ] Data looks correct
- [ ] No obvious corruption
- [ ] Timestamps preserved
- [ ] Special characters displaying correctly

---

### **Phase 6: Test Application Connection** (10 minutes)

#### **Step 6.1: Update Application Configuration**
```bash
cd /Users/ishitagupta/Documents/Personal/nutriediet-go

# Your .env should already have Aiven credentials:
cat .env | grep DB_

# Should show:
# DB_USER=avnadmin
# DB_PASSWORD=AVNS_7QDxgZDlRhQXAx3QV4z
# DB_HOST=nutriediet-mysql-ishitagupta-5564.f.aivencloud.com
# DB_PORT=22013
# DB_NAME=defaultdb
```

---

#### **Step 6.2: Test Application**
```bash
# Start application
cd /Users/ishitagupta/Documents/Personal/nutriediet-go
go run main.go

# Should see:
# ✅ Database connected successfully (Host: nutriediet-mysql-ishitagupta-5564.f.aivencloud.com, Database: defaultdb)

# In another terminal, test endpoints
curl http://localhost:8080/health

# Test login (replace with real credentials)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your_email@example.com","password":"your_password"}'
```

**✅ Testing Checklist:**
- [ ] Application starts without errors
- [ ] Database connection successful
- [ ] Can fetch data from database
- [ ] Can insert data to database
- [ ] No SQL errors in logs

---

### **Phase 7: Cleanup & Documentation** (5 minutes)

#### **Step 7.1: Archive Migration Files**
```bash
# Create archive
cd ~/Desktop
tar -czf nutriediet_migration_to_aiven_$(date +%Y%m%d).tar.gz nutriediet_migration_*

# Move to safe location
mkdir -p ~/Documents/database_backups
mv nutriediet_migration_to_aiven_*.tar.gz ~/Documents/database_backups/

# Keep source backup accessible for 30 days
```

---

#### **Step 7.2: Document Migration**
```bash
cat > ~/Documents/database_backups/migration_log.txt << EOF
=== NutrieDiet Database Migration ===
Date: $(date)
Direction: External MySQL → Aiven Cloud
Status: SUCCESS

Source Database:
- Host: YOUR_SOURCE_HOST
- Database: YOUR_SOURCE_DATABASE_NAME
- Records Migrated: [INSERT COUNT]

Target Database:
- Host: nutriediet-mysql-ishitagupta-5564.f.aivencloud.com
- Port: 22013
- Database: defaultdb
- Final Record Count: [INSERT COUNT]

Backup Location: ~/Documents/database_backups/nutriediet_migration_to_aiven_$(date +%Y%m%d).tar.gz

Notes:
- [Any special considerations]
- [Tables with issues, if any]
- [Future actions needed]
EOF

cat ~/Documents/database_backups/migration_log.txt
```

---

## 🔄 If Migration Fails - Rollback Plan

### **Restore Aiven to Previous State**
```bash
# If you backed up Aiven before migration
cd ~/Desktop/nutriediet_migration_*

gunzip aiven_before_migration.sql.gz

mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  < aiven_before_migration.sql
```

---

## 📋 Quick Migration Checklist

### **Before You Start**
- [ ] Have source database credentials
- [ ] Can connect to source database
- [ ] Have backup of source database
- [ ] Can connect to Aiven
- [ ] Have backed up current Aiven state (if needed)
- [ ] Application is stopped (or not in use)

### **During Migration**
- [ ] Source backup created successfully
- [ ] Source statistics documented
- [ ] Aiven cleaned/prepared
- [ ] Data imported to Aiven
- [ ] No import errors

### **After Migration**
- [ ] All tables present in Aiven
- [ ] Record counts match
- [ ] Sample data verified
- [ ] Application connects successfully
- [ ] Application functionality tested
- [ ] Migration documented

---

## ⚠️ Common Issues & Solutions

### **Issue 1: Character Set Problems**
**Symptom:** Special characters showing as `???` or garbled

**Solution:**
```bash
# Add charset flags to import
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  --default-character-set=utf8mb4 \
  defaultdb \
  < source_full_backup.sql
```

---

### **Issue 2: SSL Connection Errors**
**Symptom:** `ERROR: SSL connection error` or `Access denied`

**Solution:**
```bash
# Ensure --ssl-mode=REQUIRED is used
# Aiven requires SSL for all connections

# Test SSL connection
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  --ssl-verify-server-cert \
  -e "SHOW STATUS LIKE 'Ssl_cipher';"
```

---

### **Issue 3: Record Count Mismatch**
**Symptom:** Different number of records after migration

**Solution:**
```bash
# Check import log for errors
grep -i "error\|warning" import_log.txt

# Re-export and re-import specific table
mysqldump -h YOUR_SOURCE_HOST \
  -u YOUR_SOURCE_USER \
  -p'YOUR_SOURCE_PASSWORD' \
  YOUR_SOURCE_DATABASE TABLE_NAME | \
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  defaultdb
```

---

### **Issue 4: Import Timeout**
**Symptom:** Connection lost during large imports

**Solution:**
```bash
# Use table-by-table import instead (see Step 4.4)
# Or increase timeout:
mysql -h nutriediet-mysql-ishitagupta-5564.f.aivencloud.com \
  -P 22013 \
  -u avnadmin \
  -p'AVNS_7QDxgZDlRhQXAx3QV4z' \
  --ssl-mode=REQUIRED \
  --connect-timeout=3600 \
  --max_allowed_packet=1G \
  defaultdb \
  < source_full_backup.sql
```

---

## 📊 Migration Timeline

| Phase | Duration | Critical? |
|-------|----------|-----------|
| Information Gathering | 10 min | Yes |
| Backup & Analyze | 20 min | Yes |
| Prepare Target | 15 min | Yes |
| Data Migration | 30 min | Yes |
| Verification | 15 min | Yes |
| Testing | 10 min | Yes |
| Cleanup | 5 min | No |
| **TOTAL** | **~1.5 hours** | - |

---

## 🎯 Success Criteria

- [ ] All tables migrated successfully
- [ ] Record counts match (within 5%)
- [ ] Sample data looks correct
- [ ] Application connects to Aiven
- [ ] No errors in application logs
- [ ] Source database backed up
- [ ] Migration documented

---

## 📞 Need Help?

If you encounter issues:
1. Check the "Common Issues" section above
2. Review import logs for specific errors
3. Verify source database credentials
4. Ensure Aiven credentials are correct
5. Test connectivity separately

---

**Created:** October 5, 2025  
**Version:** 1.0


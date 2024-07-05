# การออกแบบและการดำเนินงานของโมเดล Multi-Tenancy

## ภาพรวมการออกแบบ

### โมเดล Multi-Tenancy

ในสถาปัตยกรรม multi-tenancy ผู้เช่าหลายราย (clients) ใช้แอปพลิเคชันและฐานข้อมูลเดียวกัน โดยข้อมูลของแต่ละผู้เช่าจะถูกแยกจากกัน ข้อมูลของผู้เช่าทั้งหมดจะถูกเก็บในสคีมาเดียวกัน

### ฐานข้อมูลเดียว สคีมาเดียว

- **ฐานข้อมูลเดียว**: ข้อมูลของผู้เช่าทั้งหมดจะถูกเก็บในฐานข้อมูลเดียว ลดภาระงานและทำให้ง่ายต่อการจัดการ และ แบ่งข้อมูลโดย tenant ID ในที่นี้จะเป็น company กับ branch
- **สคีมาเดียว**: สคีมาเดียวกันจะถูกใช้ร่วมกันในผู้เช่าทั้งหมด วิธีนี้ช่วยลดการทำสำเนาสคีมาและลดความซับซ้อนในการบำรุงรักษา

### Table Partitioning

Table Partitioning ใช้เพื่อเพิ่มประสิทธิภาพและการจัดการ ข้อมูลของผู้เช่าแต่ละรายจะถูกเก็บในพาร์ติชั่นแยกภายในตารางเดียว

#### ข้อดีของ Table Partitioning

- **ประสิทธิภาพ**: การแบ่งพาร์ติชั่นช่วยเพิ่มประสิทธิภาพการสืบค้นข้อมูลโดยการสแกนเฉพาะพาร์ติชั่นที่เกี่ยวข้อง
- **การจัดการ**: พาร์ติชั่นสามารถจัดการ เพิ่ม หรือลบได้ง่ายโดยไม่ส่งผลกระทบต่อตารางทั้งหมด
- **การขยายตัว**: สถาปัตยกรรมรองรับการขยายตัวได้อย่างมีประสิทธิภาพ โดยสามารถเพิ่มผู้เช่าใหม่ได้โดยการสร้างพาร์ติชั่นใหม่

## Installation

### Steps
1. Clone the repository:
   
```sh
 git clone https://github.com/NonPkrub/Multi-tenancy.git
 cd Multi-tenancy
```

2. Set up the database:

```sh
   CREATE SCHEME COMPANY

   \dt

   CREATE DATABASE onesystem;

```

3. Run the SQL scripts to set up tables and partitions:

```sh
    psql -U your-username -d onesystem -f schema.sql
```

4. Configure the environment variables:

```sh
cp config.yaml.example config.yaml


database:
  host: localhost
  port: 5432
  user: your_username
  password: your_password
  dbname: onesystem

```

5. Install dependencies and run the server:

```sh
go mod download

go run ./cmd/main.go
```

## การดำเนินงาน ใน Go และ SQL

### สร้าง Table หลักมาเป็น Table name

โดยในการสร้าง table หลัก ในการสร้าง partition จำเป็นต้อง partition ด้วย PRIMARY KEY นั้นหมายถึงชื่อ partition ของแต่ละ partition จะไม่ซ้ำกัน ในที่นี้จะเป็น company ที่ต้องตั้งชื่อไม่ซ้ำกัน

```sql
create table company.onesystem (
 	company varchar(255) not null,
 	branch varchar(255) not null,
 	id uuid DEFAULT gen_random_uuid()
 	first_name varchar(255) not null,
 	last_name varchar(255) not null,
 	username varchar(255) not null,
 	password varchar(255) not null,
 	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 	update_at TIMESTAMP,
 	delete_at TIMESTAMP,
 	role varchar(255) default 'user',
 	PRIMARY KEY (company, branch, uuid)
) partition by List (company);
```

### การสร้าง Company

ใน Go จะเรียก
```sh
Post("/company", s.manage.CreateCompany)
```
**_Body_** 
```sh
"company":"company_name"
```

sql query:

```sql
CREATE TABLE company.company_name PARTITION OF company.onesystem
    FOR VALUES IN ('company_name')
PARTITION BY LIST (branch)
```

### การสร้าง Branch

ใน Go จะเรียก
```sh
Post("/branch", s.manage.CreateBranch)
```
**_Body_** 
```sh
"company": "company_name",
"branch": "branch_name"
```

sql query:

```sql
CREATE TABLE company.branch_name PARTITION OF company.company_name
    FOR VALUES IN ('table_name')
```

### การเรียกดู Company

ใน Go จะเรียก
 ```sh
Get("/company", s.manage.GetCompany)
```
sql query:

```sql
SELECT inhrelid::regclass AS company
		FROM pg_inherits
		WHERE inhparent = 'company.onesystem'::regclass;
```

### การเรียกดู Branch

ใน Go จะเรียก
```sh
Get("/branch/:company", s.manage.GetBranch)
```
**_Parameter_** 
```sh
:company = company_name
```
sql query:

```sql
SELECT  inhparent::regclass AS company,inhrelid::regclass AS branch
		FROM pg_inherits
		WHERE inhparent = 'company.company_name''::regclass;
```

### การลบ Company

ใน Go จะเรียก
```sh
Delete("/company/:company", s.manage.DeleteCompany)
```

**_Body_** 
```sh
:company = company_name
```
sql query:

```sql
DROP TABLE company.company_name
```

### การลบ Branch

ใน Go จะเรียก

```sh
Delete("/company/:company/branch/:branch", s.manage.DeleteBranch)
```
**_Parameter_** 
```sh
:company = company_name , :branch = branch_name
```
ส่ง company_name เพื่อ check valid company

sql query:

```sql
DROP TABLE company.branch_name
```

### check valid company or branch

sql query:

```sql
SELECT EXISTS (
            SELECT * FROM information_schema.tables
            WHERE table_schema = 'company'
            AND table_name = company_name or branch_name
        )
```

### การสร้าง Company

ใน Go จะเรียก
```sh
Post("/company", s.manage.CreateCompany)
```
**_Body_** 
```sh
"company":"company_name"
```

sql query:

```sql
CREATE TABLE company.company_name PARTITION OF company.onesystem
    FOR VALUES IN ('company_name')
PARTITION BY LIST (branch)
```

### การอัพ Company ไปเป็น Branch

ใน Go จะเรียก
```sh
Put("/company/:company", s.manage.UpdateCompanyToBranch)
```
**_Parameter_** 
```sh
"company":"company_name(old)"
```
**_Body_**
```sh
"branch": "branch_name(old)",
"new_company": "new_company_name",
"new_branch": "new_branch_name",
"branch_name": "value_in_new_partition",
```
sql query:

```sql
-- step 1: detach partition
ALTER TABLE company.company_name(old) DETACH PARTITION company.branch_name(old);

--step 2: create branch
CREATE TABLE company.new_branch_name PARTITION OF company.new_company_name  FOR VALUES IN ('value_in_new_partition');

--step 3: insert data (copy from old table)
INSERT INTO company.new_branch_name (company,branch,id,first_name,last_name,username,password, create_at, update_at,delete_at, role)
SELECT 'new_company_name','new_branch_name',id,first_name,last_name,username,password, create_at, update_at,delete_at, role
FROM company.branch_name(old);

--step 4: update company
UPDATE company.onesystem
SET company = 'new_company_name', branch = 'new_branch_name'
WHERE company = 'company_name(old)' AND branch = 'branch_name(old)'

--step 5: delete old company
DROP TABLE company.branch_name(old)
--then
DROP TABLE company.company_name(old)

```

### การอัพ Branch ไปเป็น Company

ใน Go จะเรียก
```sh
Put("/branch/:branch", s.manage.UpdateBranchToCompany)
```
**_Parameter_** 
```sh
"branch":"branch_name(old)"
```

**_Body_**
```sh
"company": "company_name(old)",
"new_company": "new_company_name",
"new_branch": "new_branch_name",
"branch_name": "value_in_new_partition",
```
sql query:

```sql
-- step 1: detach partition
ALTER TABLE company.company_name(old) DETACH PARTITION company.branch_name(old);

--step 2: create company
CREATE TABLE company.new_company_name PARTITION OF company.onesystem  FOR VALUES IN ('new_company_name')
PARTITION BY LIST (branch);

--step 3: create branch
CREATE TABLE company.new_branch_name PARTITION OF company.new_company_name FOR VALUES IN ('value_in_new_partition');

--step 4: insert data (copy from old table)
INSERT INTO company.new_branch_name (company,branch,id,first_name,last_name,username,password, create_at, update_at,delete_at, role)
SELECT 'new_company_name','new_branch_name',id,first_name,last_name,username,password, create_at, update_at,delete_at, role
FROM company.branch_name(old);

--step 5: update company
UPDATE company.onesystem
SET company = 'new_company_name'
WHERE company = 'company_name(old)' AND branch = 'branch_name(old)'

--step 6: delete old company
DROP TABLE company.branch_name(old)

```

### การเปลี่ยนชื่อ Company

ใน Go จะเรียก
```sh
Put("/rename/company/:company", s.manage.UpdateCompanyName)
```
**_Parameter_** 
```sh
"company":"company_name(old)"
```
**_Body_** 
```sh
"new_company": "new_company_name"
```
sql query:

```sql
--step 1: rename company
ALTER TABLE company.company_name(old) RENAME TO new_company_name;

--step 2: detach company
ALTER TABLE company.onesystem DETACH PARTITION company.new_company_name;

--step 3: update company
UPDATE company.new_company_name SET branch = 'new_company_name' WHERE company = 'company_name(old)';

--step 4: attach company
ALTER TABLE company.onesystem ATTACH PARTITION company.new_company_name FOR VALUES IN ('new_company_name');

```

### การเปลี่ยนชื่อ Brand

ใน Go จะเรียก
```sh
Put("/rename/branch/:branch", s.manage.UpdateBranchName)
```
**_Parameter_** 
```sh
"branch":"branch_name(old)"
```
**_Body_** 
```sh
"company": "company_name(old)",
"new_branch": "new_branch_name"
```
sql query:
```sql
--step 1: rename branch
ALTER TABLE company.branch_name(old) RENAME TO new_branch_name;

--step 2: detach branch
ALTER TABLE company.company_name(old) DETACH PARTITION company.new_branch_name;

--step 3: update branch
UPDATE company.new_branch_name SET branch = 'new_branch_name' WHERE company = 'company_name(old)'AND branch = 'branch_name(old)';

--step 4: attach branch
ALTER TABLE company.company_name(old) ATTACH PARTITION company.new_branch_name FOR VALUES IN ('new_branch_name');
```

### Super admin

ในการออกแบบ แบบนี้ สำหรับ super admin จำเป็นต้อง initalization compnay branch และสร้างไว้ 1 record เพื่อ interaction กับ ฟังก์ชั่น การจัดการ company กับ branch

### Interaction กับข้อมูลใน onesystem Table

การสร้าง, การค้นหา, อัพเดต ,ลบ ข้อมูล สามารถ เรียก main table ใน การออก แบบระบบหลังบ้าน
ในกรณีนี้สามารถ เขียน query โดยเรียก company.onesystem ได้เลย

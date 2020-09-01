# Cloud Data Storage for MorphGNT

This project lets you import MorphGNT and WLC into AWS DynamoDB, GCP Datastore, Azure Table Storage, or any variant of SQL Server. You can also generate JSONL files for use with GCP Big Query and AWS Athena.

# Application Setup

Get Repo and get morph text

    git clone https://github.com/davidbetz/morph
    cd morph
    git clone https://github.com/morphgnt/sblgnt morphgnt

The [morphhb](https://github.com/openscriptures/morphhb) text isn't usable directly from their repo; it's in a format from the 90s. The corrected format been provided as part of this repo at `./morphwlc/hebrew` (Hebrew versification) and `./morphwlc/remapped (English versification).

Test the local setup with the following:

    make linux-print && MODE=gnt ./morph-print
    make linux-print && MODE=wlc ./morph-print

For all following examples, `MODE` can be `gnt` or `wlc`.

Use `TABLE_NAME` to explicitly set the destination. By default `morphgnt` and `morphwlc` are used.

When using `wlc`, you can use either Hebrew or English versification. Specify `VERSE_MODE=english` to use English. The default is Hebrew.

Windows is also supported:

    make windows-print
    MODE=gnt morph-print.exe

## Ephemeral VM Setup

Day to day, I install nothing. It's just Docker. When doing cloud-first, I use ephemeral VMs.

Run something at the following location:

[https://startup.netfxharmonics.com/scripts/golang/](https://startup.netfxharmonics.com/scripts/golang/)

### MSSQL

One of the benefits of cloud is that the offerings are already running. For local SQL Server, you'd need to start it. Since installings things is insane, just run it via Docker:

    docker run --name mssql --rm -dt -p 1433:1433 -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=YOUR_PASSWORD' microsoft/mssql-server-linux:2017-latest

Use [https://passwordsgenerator.net/](https://passwordsgenerator.net/) to generate a password.

Of course, any cloud offering of Microsoft SQL (e.g. Azure SQL, AWS RDS) would work too.

# Cloud Data Setup

Never create infrastructure in your application. Separate it. If your application has permissions to create resources, it has far too many permissions. It's such a poor practice that I won't even do it in a sample application.

Per the principle of least-privlege, only the listed permissions are required on the ephemeral VM.

The default table name is `morph`. You can rename this with the TABLE_NAME environment variable.

## AWS

Create the following DynamoDB table:

    Name: MorphGNT
        Hash (partition key): verse (S)
        Range (sort key): id (N)
        
    Name: MorphWLC
        Hash (partition key): verse (S)
        Range (sort key): id (N)

Client VM only requires `dynamodb:BatchWriteItem`

Run with:

    make linux-aws && MODE=gnt AWS_REGION=<REGION> ./morph-aws

You can set the region in other standard ways too.


## Azure

Create a storage account with `morphtgnt` table. Create a SAS token (with full URI).

Run with:

    make linux-azure && MODE=gnt AZURE_CS=<SAS_STRING> ./morph-azure

## GCP

Create a project with Firestore in Datastore mode

Client VM only requires `Cloud Datastore User`

Run with:

    make linux-gcp && MODE=gnt GCP_PROJECT_ID=<PROJECT_ID> ./morph-gcp

## Microsoft SQL Server

Create any database. When Azure SQL databases, you create the database and connect to the database (vs. connecting to the server).

The following will create a database via Docker (not for use with Azure SQL databases):

    docker run -it mcr.microsoft.com/mssql-tools /opt/mssql-tools/bin/sqlcmd -S $SERVER_IP -U sa -P $YOUR_PASSWORD -Q "CREATE DATABASE morph;"

SQL Server connection strings can take different forms. The following is common:

    Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;

See [https://www.connectionstrings.com/sql-server/](https://www.connectionstrings.com/sql-server/) for more details.

Run with (set MSSQL_CS):

    export MSSQL_CS='Server=SERVER_NAME;Database=morph;User Id=sa;Password=PASSWORD'
    make linux-mssql && MODE=gnt MSSQL_CS=$MSSQL_CS ./morph-mssql

Yes, `sa` is fine for local playing around. You'll be fine.

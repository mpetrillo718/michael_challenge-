# Terraform module to provision a scalable and secure static web application in AWS

# Import the AWS provider
provider "aws" {
    region = "us-east-1"
}

# Create a VPC
resource "aws_vpc" "example" {
    cidr_block = "10.0.0.0/16"
    tags = {
        Name = "example-vpc"
    }
}

# Create a subnet
resource "aws_subnet" "example" {
    vpc_id     = aws_vpc.example.id
    cidr_block = "10.0.1.0/24"
    tags = {
        Name = "example-subnet"
    }
}

# Create an Internet Gateway
resource "aws_internet_gateway" "example" {
    vpc_id = aws_vpc.example.id
    tags = {
        Name = "example-igw"
    }
}

# Create a route table
resource "aws_route_table" "example" {
    vpc_id = aws_vpc.example.id
    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.example.id
    }
    tags = {
        Name = "example-route-table"
    }
}

# Create a route table association
resource "aws_route_table_association" "example" {
    subnet_id      = aws_subnet.example.id
    route_table_id = aws_route_table.example.id
}

# Create a security group
resource "aws_security_group" "example" {
    name        = "example-security-group"
    description = "Security group for the web application"
    vpc_id      = aws_vpc.example.id

    ingress {
        from_port   = 80
        to_port     = 80
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        from_port   = 443
        to_port     = 443
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

# Create an EC2 instance
resource "aws_instance" "example" {
    ami           = "ami-0ff8a91507f77f867"
    instance_type = "t2.micro"
    subnet_id     = aws_subnet.example.id
    security_group_ids = [aws_security_group.example.id]
    tags = {
        Name = "example-ec2"
    }
}

# Create an Elastic IP
resource "aws_eip" "example" {
    vpc = true
    tags = {
        Name = "example-eip"
    }
}

# Associate the Elastic IP with the EC2 instance
resource "aws_eip_association" "example" {
    instance_id   = aws_instance.example.id
    allocation_id = aws_eip.example.id
}

# Create a self-signed SSL certificate
resource "aws_acm_certificate" "example" {
    domain_name       = "example.com"
    validation_method = "DNS"
}

# Create a CloudFront distribution
resource "aws_cloudfront_distribution" "example" {
    origin {
        domain_name = aws_instance.example.public_ip
        origin_id   = "example-origin"
    }

    enabled             = true
    is_ipv6_enabled     = true
    default_root_object = "index.html"

    default_cache_behavior {
        allowed_methods  = ["GET", "HEAD", "OPTIONS"]
        cached_methods   = ["GET", "HEAD", "OPTIONS"]
        target_origin_id = "example-origin"

        forwarded_values {
            query_string = false
            cookies {
                forward = "none"
            }
        }

        viewer_protocol_policy = "redirect-to-https"
    }

    restrictions {
        geo_restriction {
            restriction_type = "none"
        }
    }

    viewer_certificate {
        acm_certificate_arn = aws_acm_certificate.example.arn
        ssl_support_method  = "sni-only"
    }
}

# Output the CloudFront distribution domain name
output "cloudfront_domain_name" {
    value = aws_cloudfront_distribution.example.domain_nam
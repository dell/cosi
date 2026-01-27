#!/usr/bin/env perl
# Copyright © 2023-2026 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# This software contains the intellectual property of Dell Inc.
# or is licensed to Dell Inc. from third parties. Use of this software
# and the intellectual property contained therein is expressly limited to the
# terms and conditions of the License Agreement under which it is provided by or
# on behalf of Dell Inc. or its subsidiaries.

use strict;
use warnings;
use Getopt::Long;

my $filename;
my $output_file;

GetOptions(
    'f|filename=s'   => \$filename,
    'o|output=s'     => \$output_file,
) or die "Usage: perl example-config.pl -f <filename> -o <output>\n";

die "No filename provided.\n" unless defined $filename;
die "No output file provided.\n" unless defined $output_file;

my $yaml_block = '';

# Check if the filename is a valid URL
if ($filename =~ m{^https?://}) {
    # Fetch the content using curl command
    my $content = `curl -s $filename`;
    die "Failed to fetch URL: $filename\n" unless defined $content;
    open(my $fh, '<', \$content) or die "Failed to open URL content: $!\n";
    process_file($fh);
} else {
    open(my $fh, '<', $filename) or die "Failed to open file: $filename $!\n";
    process_file($fh);
}

sub process_file {
    my ($fh) = @_;

    my $found_start = 0;

    # Look for the start of the YAML block
    while (my $line = <$fh>) {
        if ($line =~ /^```yaml/) {
            # Start of YAML block found
            $found_start = 1;
            last;
        }
    }

    # Extract the YAML block
    if ($found_start) {
        while (my $line = <$fh>) {
            if ($line =~ /^```/) {
                # End of YAML block
                last;
            }
            $yaml_block .= $line;
        }
    }

    close($fh);

    if ($found_start) {
        # Write YAML block to output file
        open(my $output_fh, '>', $output_file) or die "Failed to open output file: $output_file $!\n";
        print $output_fh $yaml_block;
        close($output_fh);
        print "YAML block extracted and saved to $output_file.\n";
    } else {
        print "No YAML block found in the file.\n";
    }
}

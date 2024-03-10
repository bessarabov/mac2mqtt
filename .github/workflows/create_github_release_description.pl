#!/usr/bin/perl

use strict;
use warnings;
use feature qw(say);

sub write_file {
    my ($file_name, $content) = @_;

    open(my $fh, '>', $file_name) or die "Could not open file '$file_name' $!";

    print $fh $content;
    close $fh;
}

sub main {

    my $current_tag = $ARGV[0];

    if (not defined $current_tag) {
        die "Must run as '$0 TAG'"
    }

    my $changelog_for_this_version = '';
    my $previous_tag = '';
    my $is_in_changelog_for_this_version = 0;

    my $changelog_file_name = 'CHANGELOG.md';

    open(my $fh, '<', $changelog_file_name) or die "Could not open file '$changelog_file_name' $!";

    while (my $line = <$fh>) {
        chomp $line;

        if ($line =~ /^\Q$current_tag\E\s+/) {
            $is_in_changelog_for_this_version = 1;
        } elsif ($line =~ /^(\d+\.\d+\.\d+)\s+/a) {
            $is_in_changelog_for_this_version = 0;
            $previous_tag = $1;
            last;
        }

        if ($is_in_changelog_for_this_version) {
            $changelog_for_this_version .= $line . "\n";
        }
    }
    close $fh;

    $changelog_for_this_version =~ s/\n+$//;

    if ($changelog_for_this_version eq '') {
        die "Can't find info about release in $changelog_file_name";
    }

    if ($previous_tag eq '') {
        die "Can't find info about previous release in $changelog_file_name";
    }

    my $release_description = "
```
$changelog_for_this_version
```

Full diff: https://github.com/bessarabov/mac2qtt/compare/$previous_tag...$current_tag
";

    write_file('release_description.txt', $release_description);

}
main();

#!/bin/zsh
# Exit on error:
set -e
#
mempassdir=~/mempass
bindir=~/go/bin
pwsafeexport=$mempassdir/pwsafe.txt
encpwsafeexport=$mempassdir/pwsafe.txt.enc
encsecrets=$mempassdir/mempass-secrets.txt.enc
plainsecrets=$mempassdir/mempass-secrets.txt
translateinput=$mempassdir/pwlist.txt
translateoutput=$mempassdir/mempass.txt
# If pwSafe export file is not available in plain text:
if [[ ! -f $pwsafeexport ]]; then
  if [[ -f $encpwsafeexport ]]; then
    echo Decrypting $encpwsafeexport
    openssl enc -d -aes-256-cbc -a -in $encpwsafeexport -out $pwsafeexport
  else
    echo no pwSafe export available
    exit 1
  fi
fi
# Translate of pwSafe export to a usable file
echo Making pwSafe export $pwsafeexport usable for translate
$bindir/pwsconvert -f $pwsafeexport -o $translateinput
# Decrypt secrets file to plain text
echo Decrypting $encsecrets to $plainsecrets
openssl enc -d -aes-256-cbc -a \
   -in $encsecrets \
   -out $plainsecrets
wait
echo ----------
if [[ ! -f $plainsecrets ]]; then
	echo Trouble with decryption
	exit 1
fi
# Translate, i.e. generate result / mempass file
echo Translating $translateinput based on $plainsecrets to $translateoutput
$bindir/pwtranslate -s $plainsecrets -f $translateinput -o $translateoutput
wait
# Remove plain text secrets
echo Removing the plain text secrets file
wait
# Encrypt pwSafe export
echo Encrypting $pwsafeexport
openssl enc -aes-256-cbc -a -salt \
   -in $pwsafeexport \
   -out $encpwsafeexport
# Delete secret plain text files not longer needed
# - pwSafe export now encrypted
rm $pwsafeexport
# - Input file derived from pwSafe export not needed any more
rm $translateinput
# - secrets / config is also available as encrypted file
rm $plainsecrets
echo Completed. Check results in $mempassdir!
echo Delete $encpwsafeexport if no longer needed.


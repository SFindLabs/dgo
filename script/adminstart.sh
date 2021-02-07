./script/adminstop.sh
nohup ./dgo_admin -t admin >> log/admin.log 2>&1 &
echo "restart ok"
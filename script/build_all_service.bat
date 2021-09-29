@echo off
set service_source_root=(../apps/gate/ ../apps/global/ ../apps/logic/ ../apps/login/ ../apps/management)
set bin_dir="../bin"
set logs_dir="../logs"

IF NOT EXIST %bin_dir% (md %bin_dir%)
IF NOT EXIST %logs_dir% (md %logs_dir%)

set begin_path=%cd%

for %%A in %service_source_root% do (
cd %begin_path%
cd %%A
@echo on
make install target=windows
@echo off
)
cd %begin_path%
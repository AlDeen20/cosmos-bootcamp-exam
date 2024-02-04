#!/usr/bin/env bash

where=$(dirname "$0")

$where/x/tollroad/testing.sh $where
genesisResult=$(echo $?)
if [[ $genesisResult -eq 255 ]];
then
    exit 255
fi
genesisMaxFails=1
genesisWeight=1
genesisWins=$(echo "scale=3;((100*($genesisMaxFails-$genesisResult)*$genesisWeight)/$genesisMaxFails)" | bc)

$where/x/tollroad/roadoperatorstudent/testing.sh $where
roadOperatorResult=$(echo $?)
if [[ $roadOperatorResult -eq 255 ]];
then
    exit 255
fi
roadOperatorMaxFails=11
roadOperatorWeight=2
roadOperatorWins=$(echo "scale=3;((100*($roadOperatorMaxFails-$roadOperatorResult)*$roadOperatorWeight)/$roadOperatorMaxFails)" | bc)

$where/x/tollroad/uservaultstudent/testing.sh $where
userVaultResult=$(echo $?)
if [[ $userVaultResult -eq 255 ]];
then
    exit 255
fi
userVaultMaxFails=46
userVaultWeight=4
userVaultWins=$(echo "scale=3;((100*($userVaultMaxFails-$userVaultResult)*$userVaultWeight)/$userVaultMaxFails)" | bc)

$where/testing-cosmjs.sh
cosmjsResult=$(echo $?)
if [[ $cosmjsResult -eq 0 ]];
then
    # It's a pass
    :
elif [[ $cosmjsResult -ge 11 ]];
then
    cosmjsResult=$(($cosmjsResult-10))
else
    exit 255
fi
cosmjsMaxFails=4
cosmjsWeight=2
cosmjsWins=$(echo "scale=3;((100*($cosmjsMaxFails-$cosmjsResult)*$cosmjsWeight)/$cosmjsMaxFails)" | bc)

totalWeights=$(($genesisWeight+$roadOperatorWeight+$userVaultWeight+$cosmjsWeight))

score=$(echo "scale=3;(($genesisWins+$roadOperatorWins+$userVaultWins+$cosmjsWins)/$totalWeights)" | bc)

echo "FS_SCORE:"$score"%"
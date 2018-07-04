//生成日期
var nowTime = new Date();
var nowYear = nowTime.getFullYear();
var nowMonth = nowTime.getMonth() + 1;
var nowDate = nowTime.getDate() ;
function creatBeginDate(){
  for(var i = 2018; i<=nowYear;i++){
    var option = document.createElement('option');
    option.setAttribute('value',i);
    option.innerHTML = i;
    beginYear.appendChild(option);
  }
  //生成1月-12月
  for(var i = 1; i <=nowMonth; i++){
    var option = document.createElement('option');
    option.setAttribute('value',i);
    option.innerHTML = i;
    beginMonth.appendChild(option);
  }
  //生成1日—31日
  for(var i = 1; i <=31; i++){
    var option = document.createElement('option');
    option.setAttribute('value',i);
    option.innerHTML = i;
    beginDay.appendChild(option);
  }
}
creatBeginDate();
//保存某年某月的天数
var bdays;
//年份点击
beginYear.onclick = function(){
  //月份显示默认值
  beginMonth.options[0].selected = true;
  //天数显示默认值
  beginDay.options[0].selected = true;
};
//月份点击
$("#beginMonth").change(function(){
  //天数显示默认值
  beginDay.options[0].selected = true;
  //计算天数的显示范围
  //如果是2月
  if(beginMonth.value==nowMonth && beginYear.value==nowYear){
    for(j=bdays+1;j<=31;j++){
      beginDay.add(new Option(j,j),undefined)
    }
    for(i=31;i>=nowDate+1;i--){
      beginDay.remove(i);
    }
    return
  }
  if(beginMonth.value == 2){
    //如果是闰年
    if((beginYear.value % 4 === 0 && beginYear.value % 100 !== 0)  || beginYear.value % 400 === 0){
      bdays = 29;
      //如果是平年
    }else{
      bdays = 28;
    }
    //如果是第4、6、9、11月
  }else if(beginMonth.value == 4 || beginMonth.value == 6 ||beginMonth.value == 9 ||beginMonth.value == 11){
    bdays = 30;
  }else{
    bdays = 31;
  }
  //增加或删除天数
  //如果是28天，则删除29、30、31天(即使他们不存在也不报错)
  if(bdays == 28){
    beginDay.remove(31);
    beginDay.remove(30);
    beginDay.remove(29);
  }
  //如果是29天
  if(bdays == 29){
    beginDay.remove(31);
    beginDay.remove(30);
    //如果第29天不存在，则添加第29天
    if(!beginDay.options[29]){
      beginDay.add(new Option('29','29'),undefined)
    }
  }
  //如果是30天
  if(bdays == 30){
    beginDay.remove(31);
    //如果第29天不存在，则添加第29天
    if(!beginDay.options[29]){
      beginDay.add(new Option('29','29'),undefined)
    }
    //如果第30天不存在，则添加第30天
    if(!beginDay.options[30]){
      beginDay.add(new Option('30','30'),undefined)
    }
  }
  //如果是31天
  if(bdays == 31){
    //如果第29天不存在，则添加第29天
    if(!beginDay.options[29]){
      beginDay.add(new Option('29','29'),undefined)
    }
    //如果第30天不存在，则添加第30天
    if(!beginDay.options[30]){
      beginDay.add(new Option('30','30'),undefined)
    }
    //如果第31天不存在，则添加第31天
    if(!beginDay.options[31]){
      beginDay.add(new Option('31','31'),undefined)
    }
  }
});

//生成日期
function creatEndDate(){
  //生成1900年-2100年
  for(var i = 2018; i<=nowYear;i++){
    var option = document.createElement('option');
    option.setAttribute('value',i);
    option.innerHTML = i;
    endYear.appendChild(option);
  }
  //生成1月-12月
  for(var i = 1; i <=nowMonth; i++){
    var option = document.createElement('option');
    option.setAttribute('value',i);
    option.innerHTML = i;
    endMonth.appendChild(option);
  }
  //生成1日—31日
  for(var i = 1; i <=31; i++){
    var option = document.createElement('option');
    option.setAttribute('value',i);
    option.innerHTML = i;
    endDay.appendChild(option);

  }
}
creatEndDate();
//保存某年某月的天数
var edays;
//年份点击
endYear.onclick = function(){
  //月份显示默认值
  endMonth.options[0].selected = true;
  //天数显示默认值
  endDay.options[0].selected = true;
}
//月份点击
$("#endMonth").change(function(){
  //天数显示默认值

  endDay.options[0].selected = true;
  //计算天数的显示范围
  //如果是2月

  if(endMonth.value==nowMonth && endYear.value==nowYear){
    for(ej=edays+1;ej<=31;ej++){
      endDay.add(new Option(ej,ej),undefined)
    }
    for(ei=31;ei>=nowDate+1;ei--){
      endDay.remove(ei);
    }
    return
  }

  if(endMonth.value == 2){
    //如果是闰年
    if((endYear.value % 4 === 0 && endYear.value % 100 !== 0)  || endYear.value % 400 === 0){
      edays = 29;
      //如果是平年
    }else{
      edays = 28;
    }
    //如果是第4、6、9、11月
  }else if(endMonth.value == 4 || endMonth.value == 6 ||endMonth.value == 9 ||endMonth.value == 11){
    edays = 30;
  }else{
    edays = 31;
  }
  //增加或删除天数
  //如果是28天，则删除29、30、31天(即使他们不存在也不报错)
  if(edays == 28){
    endDay.remove(31);
    endDay.remove(30);
    endDay.remove(29);
  }
  //如果是29天
  if(edays == 29){
    endDay.remove(31);
    endDay.remove(30);
    //如果第29天不存在，则添加第29天
    if(!endDay.options[29]){
      endDay.add(new Option('29','29'),undefined)
    }
  }
  //如果是30天
  if(edays == 30){
    endDay.remove(31);
    //如果第29天不存在，则添加第29天
    if(!endDay.options[29]){
      endDay.add(new Option('29','29'),undefined)
    }
    //如果第30天不存在，则添加第30天
    if(!endDay.options[30]){
      endDay.add(new Option('30','30'),undefined)
    }
  }
  //如果是31天
  if(edays == 31){
    //如果第29天不存在，则添加第29天
    if(!endDay.options[29]){
      endDay.add(new Option('29','29'),undefined)
    }
    //如果第30天不存在，则添加第30天
    if(!endDay.options[30]){
      endDay.add(new Option('30','30'),undefined)
    }
    //如果第31天不存在，则添加第31天
    if(!endDay.options[31]){
      endDay.add(new Option('31','31'),undefined)
    }
  }
});

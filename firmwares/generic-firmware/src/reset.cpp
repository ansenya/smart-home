#include "reset.h"
#include "state.h"

void checkReset(){
  pinMode(RESET_PIN, INPUT_PULLUP);
  if(digitalRead(RESET_PIN) == LOW){
    LOG("[RESET] Button pressed");
    delay(3000);
    if(digitalRead(RESET_PIN) == LOW){
      LOG("[RESET] Factory reset!");
      prefs.clear();
      ESP.restart();
    }
  }
}

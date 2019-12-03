package snake

import (
  "math/rand"
  "time"
)
/*
func chase(a *arena){

  e  := a.enemy
  s  := a.snake

  xD := s.body[len(e.body)-1].x - e.body[len(e.body)-1].x
  yD := s.body[len(e.body)-1].y - e.body[len(e.body)-1].y
  switch e.direction{
  case RIGHT:
    if math.Abs(xD)>math.Abs(yD){
      if xD>0{
        return RIGTH
      } else{
        if yD>0{
          return UP
        } else{
          return DOWN
        }
      }
    } else{
      if yD>0{
        return UP
      } else{
        return DOWN
      }
    }
    break
  case LEFT:
    if math.Abs(xD)>math.Abs(yD){
      if xD<0{
        return LEFT
      } else{
        if yD>0{
          return UP
        } else{
          return DOWN
        }
      }
    } else{
      if yD>0{
        return UP
      } else{
        return DOWN
      }
    }
    break
  case UP:
    if math.Abs(xD)<math.Abs(yD){
      if yD>0{
        return UP
      } else{
        if xD>0{
          return RIGHT
        } else{
          return LEFT
        }
      }
    } else{
      if xD>0{
        return RIGHT
      } else{
        return LEFT
      }
    }
    break
  case DOWN:
    if math.Abs(xD)<math.Abs(yD){
      if yD<0{
        return DOWN
      } else{
        if xD>0{
          return RIGHT
        } else{
          return LEFT
        }
      }
    } else{
      if xD>0{
        return RIGHT
      } else{
        return LEFT
      }
    }
    break
  }

}*/

func enemyMove(evChan chan direction, g *Game){
  for{
    a  := g.arena
    e  := a.enemy
    eX := e.body[len(e.body)-1].x
    eY := e.body[len(e.body)-1].y
    dirChange := false
    switch e.direction{
    case RIGHT:
      if eX>=a.width-4{
        if eY<eY-a.height{
          evChan <- UP
          dirChange = true
        } else{
          evChan <- DOWN
          dirChange = true
        }
      }
      break
    case LEFT:
      if eX<=4{
        if eY<eY-a.height{
          evChan <- DOWN
          dirChange = true
        } else{
          evChan <- UP
          dirChange = true
        }
      }
      break
    case UP:
      if eY>=a.height-4{
        if eX<eX-a.width{
          evChan <- LEFT
          dirChange = true
        } else{
          evChan <- RIGHT
          dirChange = true
        }
      }
      break
    case DOWN:
      if eY<=4{
        if eX<eX-a.width{
          evChan <- RIGHT
          dirChange = true
        } else{
          evChan <- LEFT
          dirChange = true
        }
      }
      break
    }


    if(!dirChange){
      //evChan <- chase(a)
      rD := rand.Intn(4)
      dS := rand.Intn(5)
      if(dS>3){
        switch rD {
        case 0:
          evChan <- RIGHT
          break
        case 1:
          evChan <- LEFT
          break
        case 2:
          evChan <- UP
          break
        case 3:
          evChan <- DOWN
          break
        }
      }

    }




    time.Sleep(g.moveInterval())
  }
}

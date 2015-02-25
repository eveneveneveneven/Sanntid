with Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;
use  Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;

procedure exercise7 is

   Count_Failed   : exception;   -- Exception to be raised when counting fails
   Gen            : Generator;   -- Random number generator

   protected type Transaction_Manager (N : Positive) is
      entry       Finished;
      entry       Wait_Until_Aborted;
      function    Commit return Boolean;
      procedure   Signal_Abort;
   private
      Finished_Gate_Open   : Boolean := False;
      Aborted              : Boolean := False;
      Will_Commit          : Boolean := True;
   end Transaction_Manager;
   protected body Transaction_Manager is
      entry Finished when Finished_Gate_Open or Finished'Count = N is
      begin
         ------------------------------------------
         -- PART 3: Complete the exit protocol here
         ------------------------------------------
         Finished_Gate_Open := True;
        
                  
         if(Finished'Count = 0) then
            Finished_Gate_Open := False;
         end if;
         
      end Finished;
      
      entry Wait_Until_Aborted when Aborted is
      begin
         if(Wait_until_aborted'Count = 0) then
            aborted := false;
         end if;
         
      end Wait_Until_Aborted;
      
      procedure Signal_Abort is
      begin
        put_line("ABORT!!!!!");
         Aborted := True;
      end Signal_Abort;

      function Commit return Boolean is
      begin
         return Will_Commit;
      end Commit;
      
   end Transaction_Manager;

   function Unreliable_Slow_Add (x : Integer) return Integer is
   Error_Rate : Constant := 0.15;  -- (between 0 and 1)

   begin
         -------------------------------------------
      -- PART 1: Create the transaction work here
      -------------------------------------------
   
      if(Random(Gen) > Error_rate) Then
        delay Duration(Random(Gen) * 4.0);
        
        return x + 10;
      else
        delay Duration(Random(Gen) * 0.5);
        raise Count_Failed;
      end if;
   end Unreliable_Slow_Add;

   task type Transaction_Worker (Initial : Integer; Manager : access Transaction_Manager);
   task body Transaction_Worker is
      Num         : Integer   := Initial;
      Prev        : Integer   := Num;
      Round_Num   : Integer   := 0;
   begin
      Put_Line ("Worker" & Integer'Image(Initial) & " started");

      loop
         Put_Line ("Worker" & Integer'Image(Initial) & " started round" & Integer'Image(Round_Num));
         Round_Num := Round_Num + 1;

         ---------------------------------------
         -- PART 2: Do the transaction work here          
         ---------------------------------------
         
         select
            Manager.Wait_Until_Aborted;
            Num := Num + 5;
            Put_Line ("  Abort occured. Worker" & Integer'Image(Initial) & " adding 5 and comitting" & Integer'Image(Num));
            
         then abort
            
            begin
            Num := Unreliable_Slow_Add(Num);
         exception
            when Count_Failed =>
               put_line("COUNT FAILED!!!!!");
               Manager.Signal_Abort;
         end;
           Manager.Finished;
            if Manager.Commit = True then
               Put_Line ("  Worker" & Integer'Image(Initial) & " comitting" & Integer'Image(Num));
            end if;  
            
         end select;
         
         delay 0.5;

      end loop;
   end Transaction_Worker;

   Manager : aliased Transaction_Manager (3);

   Worker_1 : Transaction_Worker (0, Manager'Access);
   Worker_2 : Transaction_Worker (1, Manager'Access);
   Worker_3 : Transaction_Worker (2, Manager'Access);

begin
   Reset(Gen); -- Seed the random number generator
end exercise7;
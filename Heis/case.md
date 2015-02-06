heisnettverkstilstand = HNT

Oppstart program:
1. Listen UDP om det fins en master som broadcaster.

1a. Hvis, connect TCP som slave.
2a. Motta id fra master.

1b. Hvis ikke, bli master, begynn UDP broadcast.
2b. Send ut id fra 1 og økende til nye slaver.

=== Som master: ===========================
1. Send TCP kall 10 ganger i sekunder.
2. Over kallet send heisnettverkstilstand.
3. Kjør kostfunksjon på bestillinger
4. Ta bestilling hvis det blir den slaven.

Hvis master detter ut:
1. Slave 1 blir ny master, heretter kalt master.
2. Master sender melding til alle slaver om å dekremmentere sin id.

Hvis slave detter ut:
1. Send melding om at alle slaver med en id større enn den bortfallende dekrementerer sin id.

Bestilling:
1. Sender heisnettverkstilstand til alle slaver

3a. Hvis minst en ikke får kontakt, ikke send godkjenning.
4a. Ta ut heisene fra heisnettverkstilstand som ikke får kontakt.
5a. Husk å send ut melding om id-dekrementering.
6a. Send ut ny heisnettverkstilstand til alle resterende heiser.

3b. Hvis alle får kontakt, send godkjenning til alle heiser.

=== Som slave: ============================
1. Høre etter TCP kall fra master.
2. Motta heisnettverkstilstand over kall
3. Kjør kostfunksjon på bestillinger
4. Ta bestilling hvis det blir den slaven.

Hvis master detter ut:
1. Hvis slave har id 1, bli ny master.
2. ellers gjør ingenting

Hvis slave detter ut:
1. Vent svar fra master.

Bestilling:
1. Sender bestilling til master.
2. Får heisnettverkstilstand fra master.
3. Venter på godkjenning fra master.
4a. Hvis ikke, bruk forrige heisnettverkstilstand.
5a. Vent på neste heisnettverkstilstand av master.

4b. Hvis godkjenning, beregn neste-state til heis.
